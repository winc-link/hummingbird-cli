/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package config

import "fmt"

var LogoContent = "\n ____  ____                                     _                  __         _                 __  \n|_   ||   _|                                   (_)                [  |       (_)               |  ] \n  | |__| |  __   _   _ .--..--.   _ .--..--.   __   _ .--.   .--./)| |.--.   __   _ .--.   .--.| |  \n  |  __  | [  | | | [ `.-. .-. | [ `.-. .-. | [  | [ `.-. | / /'`\\;| '/'`\\ \\[  | [ `/'`\\]/ /'`\\' |  \n _| |  | |_ | \\_/ |, | | | | | |  | | | | | |  | |  | | | | \\ \\._//|  \\__/ | | |  | |    | \\__/  |  \n|____||____|'.__.'_/[___||__||__][___||__||__][___][___||__].',__`[__;.__.' [___][___]    '.__.;__] \n                                                           ( ( __))                                 \n" + fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, 0, 0, 33, "A CLI tool for building hummingbird driver.", 0x1B)

var (
	Version = "1.0"

	GitHubAddr              = "https://github.com/"
	GiteeAddr               = "https://gitee.com/"
	TcpProtocolDriver       = "winc-link/hummingbird-tcp-driver"
	UdpProtocolDriver       = "winc-link/hummingbird-udp-driver"
	CoapProtocolDriver      = "winc-link/hummingbird-coap-driver"
	MqttProtocolDriver      = "winc-link/hummingbird-mqtt-driver"
	HttpProtocolDriver      = "winc-link/hummingbird-http-driver"
	WebSocketProtocolDriver = "winc-link/hummingbird-websocket-driver"
	ModbusProtocolDriver    = "winc-link/hummingbird-modbus-driver"
	OpcuaProtocolDriver     = "winc-link/hummingbird-opcua-driver"
)
