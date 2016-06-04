package itunesdll

import (
    "ml/syscall"
)

//go:cgo_import_dynamic ituneslib/itunesdll.Proc_Initialize Initialize "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_FreeMemory FreeMemory "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iTunesFreeMemory iTunesFreeMemory "iTunesHelper64.dll"

//go:cgo_import_dynamic ituneslib/itunesdll.Proc_DeviceNotificationSubscribe DeviceNotificationSubscribe "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_DeviceWaitForDeviceConnectionChanged DeviceWaitForDeviceConnectionChanged "iTunesHelper64.dll"

//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceCreate iOSDeviceCreate "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceClose iOSDeviceClose "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceConnect iOSDeviceConnect "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceDisconnect iOSDeviceDisconnect "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceGetProductType iOSDeviceGetProductType "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceGetDeviceName iOSDeviceGetDeviceName "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceGetDeviceClass iOSDeviceGetDeviceClass "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceGetProductVersion iOSDeviceGetProductVersion "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceGetCpuArchitecture iOSDeviceGetCpuArchitecture "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceGetActivationState iOSDeviceGetActivationState "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceGetUniqueDeviceID iOSDeviceGetUniqueDeviceID "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceGetUniqueDeviceIDData iOSDeviceGetUniqueDeviceIDData "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceIsJailBroken iOSDeviceIsJailBroken "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_iOSDeviceAuthorizeDsids iOSDeviceAuthorizeDsids "iTunesHelper64.dll"

//go:cgo_import_dynamic ituneslib/itunesdll.Proc_SapCreateSession SapCreateSession "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_SapCloseSession SapCloseSession "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_SapCreatePrimeSignature SapCreatePrimeSignature "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_SapVerifyPrimeSignature SapVerifyPrimeSignature "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_SapExchangeData SapExchangeData "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_SapSignData SapSignData "iTunesHelper64.dll"

//go:cgo_import_dynamic ituneslib/itunesdll.Proc_KbsyncCreateSession KbsyncCreateSession "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_KbsyncValidate KbsyncValidate "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_KbsyncGetData KbsyncGetData "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_KbsyncImport KbsyncImport "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_KbsyncCloseSession KbsyncCloseSession "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_KbsyncSaveDsid KbsyncSaveDsid "iTunesHelper64.dll"

//go:cgo_import_dynamic ituneslib/itunesdll.Proc_MachineDataStartProvisioning MachineDataStartProvisioning "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_MachineDataFinishProvisioning MachineDataFinishProvisioning "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_MachineDataFree MachineDataFree "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_MachineDataClose MachineDataClose "iTunesHelper64.dll"
//go:cgo_import_dynamic ituneslib/itunesdll.Proc_MachineDataGetData MachineDataGetData "iTunesHelper64.dll"

//go:cgo_import_dynamic ituneslib/itunesdll.Proc_EncryptJsSpToken EncryptJsSpToken "iTunesHelper64.dll"
var (
    Proc_Initialize,
    Proc_FreeMemory,
    Proc_iTunesFreeMemory,

    Proc_DeviceNotificationSubscribe,
    Proc_DeviceWaitForDeviceConnectionChanged,

    Proc_iOSDeviceCreate,
    Proc_iOSDeviceClose,
    Proc_iOSDeviceConnect,
    Proc_iOSDeviceDisconnect,
    Proc_iOSDeviceGetProductType,
    Proc_iOSDeviceGetDeviceName,
    Proc_iOSDeviceGetDeviceClass,
    Proc_iOSDeviceGetProductVersion,
    Proc_iOSDeviceGetCpuArchitecture,
    Proc_iOSDeviceGetActivationState,
    Proc_iOSDeviceGetUniqueDeviceID,
    Proc_iOSDeviceGetUniqueDeviceIDData,
    Proc_iOSDeviceIsJailBroken,
    Proc_iOSDeviceAuthorizeDsids,

    Proc_SapCreateSession,
    Proc_SapCloseSession,
    Proc_SapCreatePrimeSignature,
    Proc_SapVerifyPrimeSignature,
    Proc_SapExchangeData,
    Proc_SapSignData,

    Proc_KbsyncCreateSession,
    Proc_KbsyncValidate,
    Proc_KbsyncGetData,
    Proc_KbsyncImport,
    Proc_KbsyncCloseSession,
    Proc_KbsyncSaveDsid,

    Proc_MachineDataStartProvisioning,
    Proc_MachineDataFinishProvisioning,
    Proc_MachineDataFree,
    Proc_MachineDataClose,
    Proc_MachineDataGetData syscall.Proc

    Proc_EncryptJsSpToken syscall.Proc
)
