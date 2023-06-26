package sencell 

#MsgV1: #NFCMsg | #ReportMsg | int8 

#NFCMsg: {
    info_version: uint8 & #Current_Info_Version
    device_type: #DeviceType
    // blah blah
}

#ReportMsg: {
    info_version: uint8 & #Current_Info_Version
    // blah blah
}

#Current_Info_Version: uint8 & 0x12

#DeviceType: int & 
      #sencell_lite | 
      #router | 
      #teleport | 
      #sencell_wifi_v0 | 
      #sencell_wp_v1 | 
      #extender_v1 | 
      #testjig | 
      #device_under_test

#sencell_lite: #DeviceType & 0  
#router: #DeviceType & 1   
#teleport: #DeviceType & 2
#sencell_wifi_v0: #DeviceType & 3
#sencell_wp_v1: #DeviceType & 4
#extender_v1: #DeviceType & 5,
#testjig: #DeviceType & 6
#device_under_test: #DeviceType & 7

