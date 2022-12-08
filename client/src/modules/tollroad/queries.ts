import Long from "long"
import { createPagination, createProtobufRpcClient, QueryClient } from "@cosmjs/stargate"
import { PageResponse } from "../../types/generated/cosmos/base/query/v1beta1/pagination"
import { RoadOperator } from "../../types/generated/tollroad/road_operator"
import { SystemInfo } from "../../types/generated/tollroad/system_info"
import { UserVault } from "../../types/generated/tollroad/user_vault"
import { tollroadTypes } from "../../types/tollroad/messages"
import { tollroadDefaultRegistryTypes } from "../../tollroad_signingstargateclient"
import { Message } from "protobufjs"
import { TollroadStargateClient } from "../../tollroad_stargateclient"
import { resourceLimits } from "worker_threads"
import { QueryAllRoadOperatorResponse, QueryAllUserVaultResponse, QueryClientImpl } from "../../types/generated/tollroad/query"

export interface AllRoadOperatorResponse {
    roadOperators: RoadOperator[]
    pagination?: PageResponse
}

export interface AllUserVaultResponse {
    userVaults: UserVault[]
    pagination?: PageResponse
}

export interface TollroadExtension {
    readonly tollroad: {
        readonly getSystemInfo: () => Promise<SystemInfo>

        readonly getRoadOperator: (
            index: string,
        ) => Promise<RoadOperator | undefined>
        readonly getAllRoadOperators: (
            key: Uint8Array,
            offset: Long,
            limit: Long,
            countTotal: boolean,
        ) => Promise<AllRoadOperatorResponse>

        readonly getUserVault: (
            owner: string,
            roadOperatorIndex: string,
            token: string,
        ) => Promise<UserVault | undefined>
        readonly getAllUserVaults: (
            key: Uint8Array,
            offset: Long,
            limit: Long,
            countTotal: boolean,
        ) => Promise<AllUserVaultResponse>
    }
}

export function setupTollroadExtension(base: QueryClient): TollroadExtension { 
        const rpc = createProtobufRpcClient(base);
        // Use this service to get easy typed access to query methods
        // This cannot be used for proof verification
        const queryService = new QueryClientImpl(rpc);
       
       return {
        tollroad: {
            getSystemInfo: async () =>  {
                const { SystemInfo } = await queryService.SystemInfo({})
                return SystemInfo
            },
            getRoadOperator: async(index: string) => {
                const { roadOperator } = await queryService.RoadOperator({index: index})
                return roadOperator
            },
            getAllRoadOperators: async(key: Uint8Array, offset: Long, limit: Long,countTotal: boolean): Promise<AllRoadOperatorResponse> => {
                const response: QueryAllRoadOperatorResponse = await queryService.RoadOperatorAll({
                    pagination: {
                    key: key,
                    offset:offset,
                    limit:limit,
                    countTotal:countTotal,
                    reverse: false,
                    },
                })
                return {
                    roadOperators: response.roadOperator,
                    pagination: response.pagination,
                }

                
            },
            getUserVault: async(owner: string,roadOperatorIndex: string,token: string) => {
                const { userVault} = await queryService.UserVault({owner:owner,roadOperatorIndex:roadOperatorIndex,token:token})
                return userVault

            },
            getAllUserVaults: async(key: Uint8Array,offset: Long,limit: Long,countTotal: boolean): Promise<AllUserVaultResponse> => {
                const response: QueryAllUserVaultResponse = await queryService.UserVaultAll({
                    pagination:  {
                    key:key,
                    offset: offset,
                    limit: limit,
                    countTotal: countTotal,
                    reverse:false,
                    },

                })
                return {
                    userVaults: response.userVault,
                    pagination: response.pagination,
                }

        }
       }
    }

}
