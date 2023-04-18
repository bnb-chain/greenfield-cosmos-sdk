package runtime

import (
	appv1alpha1 "github.com/cosmos/cosmos-sdk/api/cosmos/app/v1alpha1"
	autocliv1 "github.com/cosmos/cosmos-sdk/api/cosmos/autocli/v1"
	reflectionv1 "github.com/cosmos/cosmos-sdk/api/cosmos/reflection/v1"
)

func (m appModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{Query: &autocliv1.ServiceCommandDescriptor{
		Service: appv1alpha1.Query_ServiceDesc.ServiceName,
		RpcCommandOptions: []*autocliv1.RpcCommandOptions{
			{
				RpcMethod: "Config",
				Short:     "Queries the current app config",
			},
		},
		SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
			"autocli": {
				Service: autocliv1.Query_ServiceDesc.ServiceName,
				RpcCommandOptions: []*autocliv1.RpcCommandOptions{
					{
						RpcMethod: "AppOptions",
						Short:     "Queries custom autocli options",
					},
				},
			},
			"reflection": {
				Service: reflectionv1.ReflectionService_ServiceDesc.ServiceName,
				RpcCommandOptions: []*autocliv1.RpcCommandOptions{
					{
						RpcMethod: "FileDescriptors",
						Short:     "Queries the app's protobuf file descriptors",
					},
				},
			},
		},
	}}
}
