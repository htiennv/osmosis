use cosmwasm_std::{Addr, Timestamp};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use std::cmp;

use cw_storage_plus::{Item, Map};

use crate::msg::QuotaMsg;

pub const RESET_TIME_WEEKLY: u64 = 60 * 60 * 24 * 7;

#[derive(Debug, Clone)]
pub enum FlowType {
    In,
    Out,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema, Copy)]
pub struct Flow {
    pub inflow: u128,
    pub outflow: u128,
    pub period_end: Timestamp,
}

impl Flow {
    pub fn new(inflow: impl Into<u128>, outflow: impl Into<u128>, now: Timestamp) -> Self {
        Self {
            inflow: inflow.into(),
            outflow: outflow.into(),
            period_end: now.plus_seconds(RESET_TIME_WEEKLY),
        }
    }

    pub fn balance(&self) -> u128 {
        self.inflow.abs_diff(self.outflow)
    }

    pub fn is_expired(&self, now: Timestamp) -> bool {
        self.period_end < now
    }

    // Mutating methods
    pub fn expire(&mut self, now: Timestamp) {
        self.inflow = 0;
        self.outflow = 0;
        self.period_end = now.plus_seconds(RESET_TIME_WEEKLY);
    }

    pub fn add_flow(&mut self, direction: FlowType, value: u128) {
        match direction {
            FlowType::In => self.inflow = self.inflow.saturating_add(value),
            FlowType::Out => self.outflow = self.outflow.saturating_add(value),
        }
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct Quota {
    name: String,
    max_percentage_send: u32,
    max_percentage_recv: u32,
    duration: cw_utils::Duration,
}

impl Quota {
    /// Calculates the max capacity based on the total value of the channel
    pub fn capacity_at(&self, total_value: &u128, direction: &FlowType) -> u128 {
        let max_percentage = match direction {
            FlowType::In => self.max_percentage_recv,
            FlowType::Out => self.max_percentage_send,
        };
        total_value * (max_percentage as u128) / 100_u128
    }
}

impl From<QuotaMsg> for Quota {
    fn from(msg: QuotaMsg) -> Self {
        let send_recv = (
            cmp::min(msg.send_recv.0, 100),
            cmp::min(msg.send_recv.1, 100),
        );
        Quota {
            name: msg.name,
            max_percentage_send: send_recv.0,
            max_percentage_recv: send_recv.1,
            duration: msg.duration,
        }
    }
}

/// Only this module can manage the contract
pub const GOVMODULE: Item<Addr> = Item::new("gov_module");
/// Only this module can execute transfers
pub const IBCMODULE: Item<Addr> = Item::new("ibc_module");
// For simplicity, the map keys (ibc channel) refers to the "host" channel on the
// osmosis side. This means that on PacketSend it will refer to the source
// channel while on PacketRecv it refers to the destination channel.
//
// It is the responsibility of the go module to pass the appropriate channel
// when sending the messages
pub const QUOTAS: Map<String, Vec<Quota>> = Map::new("quotas");
pub const FLOW: Map<String, Flow> = Map::new("flow");
