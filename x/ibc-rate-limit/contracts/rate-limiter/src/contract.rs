#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{Addr, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult, Timestamp};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg};
use crate::state::{Flow, FlowType, Quota, FLOW, GOVMODULE, IBCMODULE, QUOTAS};

// version info for migration info
const CONTRACT_NAME: &str = "crates.io:rate-limiter";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    env: Env,
    _info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    IBCMODULE.save(deps.storage, &msg.ibc_module)?;
    GOVMODULE.save(deps.storage, &msg.gov_module)?;

    for (channel, quotas) in msg.channel_quotas {
        QUOTAS.save(deps.storage, channel.clone(), &vec![quotas.into()])?;
        FLOW.save(
            deps.storage,
            channel,
            &Flow::new(0_u128, 0_u128, env.block.time),
        )?;
    }

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("ibc_module", msg.ibc_module.to_string())
        .add_attribute("gov_module", msg.gov_module.to_string()))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::SendPacket {
            channel_id,
            channel_value,
            funds,
        } => try_transfer(
            deps,
            info.sender,
            channel_id,
            channel_value,
            funds,
            FlowType::Out,
            env.block.time,
        ),
        ExecuteMsg::RecvPacket {
            channel_id,
            channel_value,
            funds,
        } => try_transfer(
            deps,
            info.sender,
            channel_id,
            channel_value,
            funds,
            FlowType::In,
            env.block.time,
        ),
        ExecuteMsg::AddChannel {} => todo!(),
        ExecuteMsg::RemoveChannel {} => todo!(),
    }
}

fn check_quota(
    quota: &Quota,
    flow: &mut Flow,
    direction: FlowType,
    channel_id: &str,
    channel_value: u128,
    funds: u128,
    now: Timestamp,
) -> Result<(u128, u128, Timestamp), ContractError> {
    let max = quota.capacity_at(&channel_value, &direction);
    if flow.is_expired(now) {
        flow.expire(now)
    }
    flow.add_flow(direction, funds);

    let balance = flow.balance();
    if balance > max {
        return Err(ContractError::RateLimitExceded {
            channel: channel_id.to_string(),
            reset: flow.period_end,
        });
    }
    return Ok((balance, max, flow.period_end));
}

pub fn try_transfer(
    deps: DepsMut,
    sender: Addr,
    channel_id: String,
    channel_value: u128,
    funds: u128,
    direction: FlowType,
    now: Timestamp,
) -> Result<Response, ContractError> {
    // Only the IBCMODULE can execute transfers
    let ibc_module = IBCMODULE.load(deps.storage)?;
    if sender != ibc_module {
        return Err(ContractError::Unauthorized {});
    }
    let quotas = QUOTAS.load(deps.storage, channel_id.clone())?;
    if quotas.len() == 0 {
        // No Quota configured for the current channel. Allowing all messages.
        return Ok(Response::new()
            .add_attribute("method", "try_transfer")
            .add_attribute("channel_id", channel_id)
            .add_attribute("quota", "none"));
    }

    let mut flow = FLOW.load(deps.storage, channel_id.clone())?;

    let quotas: Result<Vec<(u128, u128, Timestamp)>, _> = quotas
        .iter()
        .map(|quota| {
            check_quota(
                &quota,
                &mut flow,
                direction.clone(),
                &channel_id,
                channel_value,
                funds,
                now,
            )
        })
        .collect();
    let quotas = quotas?;

    FLOW.update(
        deps.storage,
        channel_id.clone(),
        |_| -> Result<_, ContractError> { Ok(flow) },
    )?;

    let response = Response::new()
        .add_attribute("method", "try_transfer")
        .add_attribute("channel_id", channel_id);

    // Adding the attributes from each quota to the response
    quotas.iter().fold(Ok(response), |acc, quota| {
        Ok(acc?
            .add_attribute("used", quota.0.to_string())
            .add_attribute("max", quota.1.to_string())
            .add_attribute("period_end", quota.2.to_string()))
    })

    // Ok(response
    //     .add_attribute("used", balance.to_string())
    //     .add_attribute("max", max.to_string())
    //     .add_attribute("period_end", flow.period_end.nanos().to_string()))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(_deps: Deps, _env: Env, _msg: ExecuteMsg) -> StdResult<Binary> {
    todo!()
}

#[cfg(test)]
mod tests {
    use crate::msg::QuotaMsg;
    use crate::state::RESET_TIME_WEEKLY;

    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{Addr, Attribute};

    const IBC_ADDR: &str = "IBC_MODULE";
    const GOV_ADDR: &str = "GOV_MODULE";

    #[test]
    fn proper_instantiation() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            gov_module: Addr::unchecked(GOV_ADDR),
            ibc_module: Addr::unchecked(IBC_ADDR),
            channel_quotas: vec![],
        };
        let info = mock_info(IBC_ADDR, &vec![]);

        // we can just call .unwrap() to assert this was a success
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // TODO: Check initialization values are correct
    }

    #[test]
    fn permissions() {
        let mut deps = mock_dependencies();

        let quota = QuotaMsg::new("Weekly", RESET_TIME_WEEKLY, 10, 10);
        let msg = InstantiateMsg {
            gov_module: Addr::unchecked(GOV_ADDR),
            ibc_module: Addr::unchecked(IBC_ADDR),
            channel_quotas: vec![("channel".to_string(), quota)],
        };
        let info = mock_info(IBC_ADDR, &vec![]);
        instantiate(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();

        let msg = ExecuteMsg::SendPacket {
            channel_id: "channel".to_string(),
            channel_value: 3_000,
            funds: 300,
        };

        // This succeeds
        execute(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();

        let info = mock_info("SomeoneElse", &vec![]);

        let msg = ExecuteMsg::SendPacket {
            channel_id: "channel".to_string(),
            channel_value: 3_000,
            funds: 300,
        };
        let err = execute(deps.as_mut(), mock_env(), info.clone(), msg).unwrap_err();
        assert!(matches!(err, ContractError::Unauthorized { .. }));
    }

    #[test]
    fn consume_allowance() {
        let mut deps = mock_dependencies();

        let quota = QuotaMsg::new("Weekly", RESET_TIME_WEEKLY, 10, 10);
        let msg = InstantiateMsg {
            gov_module: Addr::unchecked(GOV_ADDR),
            ibc_module: Addr::unchecked(IBC_ADDR),
            channel_quotas: vec![("channel".to_string(), quota)],
        };
        let info = mock_info(GOV_ADDR, &vec![]);
        let _res = instantiate(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();

        let msg = ExecuteMsg::SendPacket {
            channel_id: "channel".to_string(),
            channel_value: 3_000,
            funds: 300,
        };
        let info = mock_info(IBC_ADDR, &vec![]);
        let res = execute(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();
        let Attribute { key, value } = &res.attributes[2];
        assert_eq!(key, "used");
        assert_eq!(value, "300");

        let msg = ExecuteMsg::SendPacket {
            channel_id: "channel".to_string(),
            channel_value: 3_000,
            funds: 300,
        };
        let err = execute(deps.as_mut(), mock_env(), info, msg).unwrap_err();
        assert!(matches!(err, ContractError::RateLimitExceded { .. }));
    }

    #[test]
    fn symetric_flows_dont_consume_allowance() {
        let mut deps = mock_dependencies();

        let quota = QuotaMsg::new("Weekly", RESET_TIME_WEEKLY, 10, 10);
        let msg = InstantiateMsg {
            gov_module: Addr::unchecked(GOV_ADDR),
            ibc_module: Addr::unchecked(IBC_ADDR),
            channel_quotas: vec![("channel".to_string(), quota)],
        };
        let info = mock_info(GOV_ADDR, &vec![]);
        let _res = instantiate(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();

        let info = mock_info(IBC_ADDR, &vec![]);
        let send_msg = ExecuteMsg::SendPacket {
            channel_id: "channel".to_string(),
            channel_value: 3_000,
            funds: 300,
        };
        let recv_msg = ExecuteMsg::RecvPacket {
            channel_id: "channel".to_string(),
            channel_value: 3_000,
            funds: 300,
        };

        let res = execute(deps.as_mut(), mock_env(), info.clone(), send_msg.clone()).unwrap();
        let Attribute { key, value } = &res.attributes[2];
        assert_eq!(key, "used");
        assert_eq!(value, "300");

        let res = execute(deps.as_mut(), mock_env(), info.clone(), recv_msg.clone()).unwrap();
        let Attribute { key, value } = &res.attributes[2];
        assert_eq!(key, "used");
        assert_eq!(value, "0");

        // We can still use the channel. Even if we have sent more than the
        // allowance through the channel (900 > 3000*.1), the current "balance"
        // of inflow vs outflow is still lower than the channel's capacity/quota
        let res = execute(deps.as_mut(), mock_env(), info.clone(), recv_msg.clone()).unwrap();
        let Attribute { key, value } = &res.attributes[2];
        assert_eq!(key, "used");
        assert_eq!(value, "300");

        let err = execute(deps.as_mut(), mock_env(), info.clone(), recv_msg.clone()).unwrap_err();

        assert!(matches!(err, ContractError::RateLimitExceded { .. }));
        //assert_eq!(18, value.count);
    }

    #[test]
    fn asymetric_quotas() {
        let mut deps = mock_dependencies();

        let quota = QuotaMsg::new("Weekly", RESET_TIME_WEEKLY, 10, 1);
        let msg = InstantiateMsg {
            gov_module: Addr::unchecked(GOV_ADDR),
            ibc_module: Addr::unchecked(IBC_ADDR),
            channel_quotas: vec![("channel".to_string(), quota)],
        };
        let info = mock_info(GOV_ADDR, &vec![]);
        let _res = instantiate(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();

        // Sending 2%
        let msg = ExecuteMsg::SendPacket {
            channel_id: "channel".to_string(),
            channel_value: 3_000,
            funds: 60,
        };
        let info = mock_info(IBC_ADDR, &vec![]);
        let res = execute(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();
        let Attribute { key, value } = &res.attributes[2];
        assert_eq!(key, "used");
        assert_eq!(value, "60");

        // Sending 1% more. Allowed, as sending has a 10% allowance
        let msg = ExecuteMsg::SendPacket {
            channel_id: "channel".to_string(),
            channel_value: 3_000,
            funds: 30,
        };

        let info = mock_info(IBC_ADDR, &vec![]);
        let res = execute(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();
        let Attribute { key, value } = &res.attributes[2];
        assert_eq!(key, "used");
        assert_eq!(value, "90");

        // Receiving 1% should fail. 3% already executed through the channel
        let recv_msg = ExecuteMsg::RecvPacket {
            channel_id: "channel".to_string(),
            channel_value: 3_000,
            funds: 30,
        };

        let err = execute(deps.as_mut(), mock_env(), info.clone(), recv_msg.clone()).unwrap_err();
        assert!(matches!(err, ContractError::RateLimitExceded { .. }));
    }
}
