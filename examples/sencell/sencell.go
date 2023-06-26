package sencell

type MsgV1 struct {
	MsgType   MsgV1Type
	NFCMsg    *NFCMsg
	ReportMsg *ReportMsg
	Int8      int8
}

type MsgV1Type int16

const (
	NFCMsgType    MsgV1Type = 1
	ReportMsgType MsgV1Type = 2
	Int8Type      MsgV1Type = 3
)

type NFCMsg struct {
	InfoVersion uint8
	DeviceType  DeviceType
}

type ReportMsg struct {
	InfoVersion uint8
}

const CurrentInfoVersion int8 = 0x12

type DeviceType int

const (
	SencellLite     DeviceType = 0
	Router          DeviceType = 1
	Teleport        DeviceType = 2
	SencellWifiV0   DeviceType = 3
	SencellWpV1     DeviceType = 4
	ExtenderV1      DeviceType = 5
	Testjig         DeviceType = 6
	DeviceUnderTest DeviceType = 7
)
