package usecase_test

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v21/ingest/sqs/domain"
	"github.com/osmosis-labs/osmosis/v21/ingest/sqs/domain/mocks"
	"github.com/osmosis-labs/osmosis/v21/ingest/sqs/log"
	"github.com/osmosis-labs/osmosis/v21/ingest/sqs/router/usecase"
	"github.com/osmosis-labs/osmosis/v21/ingest/sqs/router/usecase/route"
)

// Tests the call to handleRoutes by mocking the router repository and pools use case
// with relevant data.
func (s *RouterTestSuite) TestHandleRoutes() {
	const (
		defaultTimeoutDuration = time.Second * 10

		tokenInDenom  = "uosmo"
		tokenOutDenom = "uion"

		minOsmoLiquidity = 10000 * usecase.OsmoPrecisionMultiplier
	)

	// Create test balancer pool

	balancerCoins := sdk.NewCoins(
		sdk.NewCoin(tokenInDenom, sdk.NewInt(1000000000000000000)),
		sdk.NewCoin(tokenOutDenom, sdk.NewInt(1000000000000000000)),
	)

	balancerPoolID := s.PrepareBalancerPoolWithCoins(balancerCoins...)
	balancerPool, err := s.App.PoolManagerKeeper.GetPool(s.Ctx, balancerPoolID)
	s.Require().NoError(err)

	defaultPool := &domain.PoolWrapper{
		ChainModel: balancerPool,
		SQSModel: domain.SQSPool{
			TotalValueLockedUSDC: osmomath.NewInt(int64(minOsmoLiquidity + 1)),
			PoolDenoms:           []string{tokenInDenom, tokenOutDenom},
			Balances:             balancerCoins,
			SpreadFactor:         DefaultSpreadFactor,
		},
	}

	var (
		defaultRoute = WithCandidateRoutePools(
			EmptyCandidateRoute,
			[]route.CandidatePool{
				{
					ID:            defaultPool.GetId(),
					TokenOutDenom: tokenOutDenom,
				},
			},
		)

		defaultSinglePools = []domain.PoolI{defaultPool}

		singleDefaultRoutes = route.CandidateRoutes{
			Routes: []route.CandidateRoute{defaultRoute},
			UniquePoolIDs: map[uint64]struct{}{
				defaultPool.GetId(): {},
			},
		}

		emptyPools = []domain.PoolI{}

		emptyRoutes = route.CandidateRoutes{}

		defaultRouterConfig = domain.RouterConfig{
			// Only these config values are relevant for this test
			// for searching for routes when none were present in cache.
			MaxPoolsPerRoute: 4,
			MaxRoutes:        4,

			// These configs are not relevant for this test.
			PreferredPoolIDs:          []uint64{},
			MaxSplitIterations:        10,
			MinOSMOLiquidity:          minOsmoLiquidity,
			RouteUpdateHeightInterval: 10,
		}
	)

	testCases := []struct {
		name string

		repositoryRoutes route.CandidateRoutes
		repositoryPools  []domain.PoolI
		takerFeeMap      domain.TakerFeeMap
		isCacheDisabled  bool

		expectedCandidateRoutes route.CandidateRoutes

		expectedError error
	}{
		{
			name: "routes in cache -> use them",

			repositoryRoutes: singleDefaultRoutes,
			repositoryPools:  emptyPools,

			expectedCandidateRoutes: singleDefaultRoutes,
		},
		{
			name: "cache is disabled in config -> recomputes routes despite having available in cache",

			repositoryRoutes: singleDefaultRoutes,
			repositoryPools:  emptyPools,
			isCacheDisabled:  true,

			expectedCandidateRoutes: emptyRoutes,
		},
		{
			name: "no routes in cache but relevant pools in store -> recomputes routes",

			repositoryRoutes: emptyRoutes,
			repositoryPools:  defaultSinglePools,

			expectedCandidateRoutes: singleDefaultRoutes,
		},
		{
			name: "no routes in cache and no relevant pools in store -> returns no routes",

			repositoryRoutes: emptyRoutes,
			repositoryPools:  emptyPools,

			expectedCandidateRoutes: emptyRoutes,
		},

		// TODO:
		// routes in cache but pools have more optimal -> cache is still used
		// multiple routes in cache -> use them
		// multiple rotues in pools -> use them
		// error in repository -> return error
		// error in storing routes after recomputing -> return error
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {

			routerRepositoryMock := &mocks.RedisRouterRepositoryMock{
				Routes: map[domain.DenomPair]route.CandidateRoutes{
					// These are the routes that are stored in cache and returned by the call to GetRoutes.
					{Denom0: tokenOutDenom, Denom1: tokenInDenom}: tc.repositoryRoutes,
				},

				// No need to set taker fees on the mock since they are only relevant when
				// set on the router for this test.
			}

			poolsUseCaseMock := &mocks.PoolsUsecaseMock{
				// These are the pools returned by the call to GetAllPools
				Pools: tc.repositoryPools,
			}

			routerUseCase := usecase.NewRouterUsecase(defaultTimeoutDuration, routerRepositoryMock, poolsUseCaseMock, domain.RouterConfig{
				RouteCacheEnabled: !tc.isCacheDisabled,
			}, &log.NoOpLogger{})

			routerUseCaseImpl, ok := routerUseCase.(*usecase.RouterUseCaseImpl)
			s.Require().True(ok)

			// Initialize router
			router := usecase.NewRouter(defaultRouterConfig.PreferredPoolIDs, defaultRouterConfig.MaxPoolsPerRoute, defaultRouterConfig.MaxRoutes, defaultRouterConfig.MaxSplitRoutes, defaultRouterConfig.MaxSplitIterations, defaultRouterConfig.MaxSplitIterations, &log.NoOpLogger{})
			router = usecase.WithSortedPools(router, poolsUseCaseMock.Pools)

			// System under test
			ctx := context.Background()
			actualCandidateRoutes, err := routerUseCaseImpl.HandleRoutes(ctx, router, tokenInDenom, tokenOutDenom)

			if tc.expectedError != nil {
				s.Require().EqualError(err, tc.expectedError.Error())
				s.Require().Len(actualCandidateRoutes, 0)
				return
			}

			s.Require().NoError(err)

			// Pre-set routes should be returned.

			s.Require().Equal(len(tc.expectedCandidateRoutes.Routes), len(actualCandidateRoutes.Routes))
			for i, route := range actualCandidateRoutes.Routes {
				s.Require().Equal(tc.expectedCandidateRoutes.Routes[i], route)
			}

			// For the case where the cache is disabled, the expected routes in cache
			// will be the same as the original routes in the repository.
			if tc.isCacheDisabled {
				tc.expectedCandidateRoutes = tc.repositoryRoutes
			}

			// Check that router repository was updated
			s.Require().Equal(tc.expectedCandidateRoutes, routerRepositoryMock.Routes[domain.DenomPair{Denom0: tokenOutDenom, Denom1: tokenInDenom}])
		})
	}
}

// Tests that routes that overlap in pools IDs get filtered out.
// Tests that the order of the routes is in decreasing priority.
// That is, if routes A and B overlap where A comes before B, then B is filtered out.
// Additionally, tests that overlapping within the same route has no effect on filtering.
// Lastly, validates that if a route overlaps with subsequent routes in the list but gets filtered out,
// then subesequent routes are not affected by filtering.
func (s *RouterTestSuite) TestFilterDuplicatePoolIDRoutes() {
	var (
		deafaultPool = &mocks.MockRoutablePool{ID: defaultPoolID}

		otherPool = &mocks.MockRoutablePool{ID: defaultPoolID + 1}

		defaultSingleRoute = WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
			deafaultPool,
		})
	)

	tests := map[string]struct {
		routes []route.RouteImpl

		expectedRoutes []route.RouteImpl
	}{
		"empty routes": {
			routes:         []route.RouteImpl{},
			expectedRoutes: []route.RouteImpl{},
		},

		"single route single pool": {
			routes: []route.RouteImpl{
				defaultSingleRoute,
			},

			expectedRoutes: []route.RouteImpl{
				defaultSingleRoute,
			},
		},

		"single route two different pools": {
			routes: []route.RouteImpl{
				WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
					deafaultPool,
					otherPool,
				}),
			},

			expectedRoutes: []route.RouteImpl{
				WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
					deafaultPool,
					otherPool,
				}),
			},
		},

		// Note that filtering only happens if pool ID duplicated across different routes.
		// Duplicate pool IDs within the same route are filtered out at a different step
		// in the router logic.
		"single route two same pools (have no effect on filtering)": {
			routes: []route.RouteImpl{
				WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
					deafaultPool,
					deafaultPool,
				}),
			},

			expectedRoutes: []route.RouteImpl{
				WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
					deafaultPool,
					deafaultPool,
				}),
			},
		},

		"two single hop routes and no duplicates": {
			routes: []route.RouteImpl{
				defaultSingleRoute,

				WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
					otherPool,
				}),
			},

			expectedRoutes: []route.RouteImpl{
				defaultSingleRoute,

				WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
					otherPool,
				}),
			},
		},

		"two single hop routes with duplicates (second filtered)": {
			routes: []route.RouteImpl{
				defaultSingleRoute,

				defaultSingleRoute,
			},

			expectedRoutes: []route.RouteImpl{
				defaultSingleRoute,
			},
		},

		"three route. first and second overlap. second and third overlap. second is filtered out but not third": {
			routes: []route.RouteImpl{
				defaultSingleRoute,

				WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
					deafaultPool, // first and second overlap
					otherPool,    // second and third overlap
				}),

				WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
					otherPool,
				}),
			},

			expectedRoutes: []route.RouteImpl{
				defaultSingleRoute,

				WithRoutePools(route.RouteImpl{}, []domain.RoutablePool{
					otherPool,
				}),
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		s.Run(name, func() {

			actualRoutes := usecase.FilterDuplicatePoolIDRoutes(tc.routes)

			s.Require().Equal(len(tc.expectedRoutes), len(actualRoutes))
		})
	}
}
