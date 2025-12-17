package client

import (
	infoapi "github.com/localstack/localstack-sdk-go/internal/generated/infoapi"
	internalapi "github.com/localstack/localstack-sdk-go/internal/generated/internalapi"
)

type Diagnose = infoapi.DiagnoseResponse

type DiagnoseVersion = infoapi.DiagnoseVersion

type Health = infoapi.HealthResponse

type Info = infoapi.InfoResponse

type LicenseInfo = infoapi.LicenseInfoResponse

type CertificateInfo = internalapi.CertificateInfo

type CertificateList = internalapi.CertificateList

type InitScriptInfo = internalapi.InitScriptInfo

type InitScripts = internalapi.InitScriptsResponse

type Plugins = internalapi.PluginsResponse

type RootCAInfo = internalapi.RootCAInfo
