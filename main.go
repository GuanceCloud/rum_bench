package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DKHost   = "localhost"
	Duration = time.Second * 10
)

const (
	AppID      = "web_abcdefg123456789"
	AppEnv     = "production"
	AppVersion = "1.0.0"
)

type RUMType string

const (
	RUMSession       RUMType = "session"
	RUMView          RUMType = "view"
	RUMResource      RUMType = "resource"
	RUMAction        RUMType = "action"
	RUMLongTask      RUMType = "long_task"
	RUMError         RUMType = "error"
	RUMSessionReplay RUMType = "session_replay"
)

var (
	rumEndpoint    = fmt.Sprintf("http://%s:9529%s", DKHost, "/v1/write/rum?precision=ns")
	replayEndPoint = fmt.Sprintf("http://%s:9529%s", DKHost, "/v1/write/rum/replay?precision=ns")
)

var (
	flagChan   = make(chan struct{}, 1)
	outputChan = make(chan *counter, 10)
)

var concurrentCnt = 1

var errorBody = `error,sdk_name=df_android_rum_sdk,sdk_version=2.0.26,app_id=___APPID___,env=___ENV___,service=browser,version=___VERSION___,userid=6931fa8d-769b-46ef-998f-34fc20947562,session_id=49c6f9c2-318b-4c99-957d-b514db8706ee,session_type=user,is_signin=F,os=Mac\\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=103.0.0.0,browser_version_major=103,screen_size=1440*900,network_type=4g,view_id=8a5fcfb3-9f18-465e-8598-5adcc9924abc,view_url=http://localhost:8080/index.html?1111,view_host=localhost:8080,view_path=/index.html,view_path_group=/index.html,view_url_query={},error_source=source,error_type=java_crash,error_handling=unhandled error_message="manual err",error_stack="java.lang.ArithmeticException: divide by zero
 at prof.wang.activity.TeamInvitationActivity.o0(Unknown Source:1)
 at prof.wang.activity.TeamInvitationActivity.k0(Unknown Source:0)
 at j9.f7.run(Unknown Source:0)
 at java.lang.Thread.run(Thread.java:1012)" ___TIMESTAMP___
error,error_source=logger,view_id=5387699e-b62f-45d5-a2b7-b70840326e1a,is_signin=T,session_id=C34678___APPID___,errer=root,error_situation=run,locale=en_US,userid=brandon.test.userid,carrier=T-Mobile,action_id=0497a099-136c-4a1d-a06e-49e5a2ce8588,action_name=sourcemap.test,error_type=ios_crash,network_type=3G,memory_total=1.92GB,view_name=SomeView,os=iOS,application_uuid=992b1f06-8fab-407c-8a47-10fc548ab7a1,os_version=15.5,custom_keys=["custom_tag"\,"track_id"],sdk_name=df_ios_rum_sdk,env=___ENV___,version=___VERSION___,device_uuid=670bfb704e385be0,sdk_package_track=1.1.0-alpha01,os_version_major=15,custom_tag=any\ tags,sdk_package_agent=1.3.6-beta02,track_id=rtrace_ded7c278d095440ca1a0595dc40e8b19,sdk_version=1.3.6-beta02,model=iPhone\ X,screen_size=1080*1794,session_type=user,app_id=___APPID___,device=APPLE error_message="divide by zero",cpu_use=1.0,error_stack="Hardware Model:  iPhone10,3
               OS Version:   iPhone OS 15.5
               Report Version:  104

               Code Type:   ARM64

               Last Exception Backtrace:
               0   CoreFoundation                      0x0000000180a8fd20 0x1809fd000 + 601376
               1   libobjc.A.dylib                     0x0000000198280ee4 objc_exception_throw + 56
               2   CoreFoundation                      0x0000000180b60750 0x1809fd000 + 1455952
               3   CoreFoundation                      0x0000000180a2960c 0x1809fd000 + 181772
               4   CoreFoundation                      0x0000000180a2881c _CF_forwarding_prep_0 + 92
               5   App                                 0x000000010022eb60 __35-[Crasher throwUncaughtNSException]_block_invoke_2 + 72
               6   AFNetworking                        0x0000000100733458 __124-[AFHTTPSessionManager dataTaskWithHTTPMethod:URLString:parameters:headers:uploadProgress:downloadProgress:success:failure:]_block_invoke_2 + 132
               7   AFNetworking                        0x0000000100744f14 __72-[AFURLSessionManagerTaskDelegate URLSession:task:didCompleteWithError:]_block_invoke_2.108 + 148
               8   libdispatch.dylib                   0x000000018074f094 0x1806ec000 + 405652
               9   libdispatch.dylib                   0x0000000180750094 0x1806ec000 + 409748
               10  libdispatch.dylib                   0x00000001807318a8 0x1806ec000 + 284840
               11  libdispatch.dylib                   0x00000001807314d0 0x1806ec000 + 283856
               12  CoreFoundation                      0x0000000180a4b0c4 0x1809fd000 + 319684
               13  CoreFoundation                      0x0000000180a085e8 0x1809fd000 + 46568
               14  CoreFoundation                      0x0000000180a1b240 CFRunLoopRunSpecific + 572
               15  GraphicsServices                    0x00000001a151d988 GSEventRunModal + 160
               16  UIKitCore                           0x000000018321b41c 0x182d36000 + 5133340
               17  UIKitCore                           0x0000000182fb4b88 UIApplicationMain + 336
               18  App                                 0x00000001002315b0 main + 96
               19  dyld                                0x00000001004a43d0 start + 444

               Binary Images:
                      0x100224000 -        0x100287fff App arm64 <5dd80683dc60325cb8bc1643733b6eab> /private/var/containers/Bundle/Application/7E63CAC0-7413-48AC-9411-8DD7FB556055/App.app/App
                      0x10072c000 -        0x100757fff AFNetworking arm64 <8177fa3cc4463476b0e946303a66b4c1> /private/var/containers/Bundle/Application/7E63CAC0-7413-48AC-9411-8DD7FB556055/App.app/Frameworks/AFNetworking.framework/AFNetworking
                      0x19826c000 -        0x1982a3fff libobjc.A.dylib arm64 <57ca31b758ea36d6a442728888f336ec> /usr/lib/libobjc.A.dylib
                      0x1809fd000 -        0x180e3afff CoreFoundation arm64 <48cd0a807a9234ebb1408c475d135808> /System/Library/Frameworks/CoreFoundation.framework/CoreFoundation
                      0x1806ec000 -        0x18076efff libdispatch.dylib arm64 <f87efead673b3a09b7ab9c69b55a18a8> /usr/lib/system/libdispatch.dylib
                      0x182d36000 -        0x1844cffff UIKitCore arm64 <a84e395e0b003162b1fdb861ce41349c> /System/Library/PrivateFrameworks/UIKitCore.framework/UIKitCore
                      0x1a151c000 -        0x1a1524fff GraphicsServices arm64 <996d6fdae7883abeb6d6ad8e0f4cc881> /System/Library/PrivateFrameworks/GraphicsServices.framework/GraphicsServices
               ",battery_use=0.0,memory_use=63.79 ___TIMESTAMP___
error,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},error_id=3f1a843d-edc2-47d2-b5a5-7c4fca855ba7,error_source=source,error_type=Error,error_handling=unhandled view_in_foreground=true,error_message="manual err",error_stack="Error: manual err
  at n @ http://127.0.0.1:8081/dist/bundle.js:2:141
  at n @ http://127.0.0.1:8081/dist/bundle.js:2:186
  at n @ http://127.0.0.1:8081/dist/bundle.js:2:179
  at n @ http://127.0.0.1:8081/dist/bundle.js:2:179
  at n @ http://127.0.0.1:8081/dist/bundle.js:2:179
  at t @ http://127.0.0.1:8081/dist/bundle.js:2:107
  at <anonymous> @ http://127.0.0.1:8081/dist/bundle.js:2:705
  at <anonymous> @ http://127.0.0.1:8081/dist/bundle.js:2:709" ___TIMESTAMP___
`

var viewBody = `view,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},view_loading_type=initial_load is_active=true,view_error_count=1,view_resource_count=10,view_long_task_count=0,view_action_count=0,first_contentful_paint=2993400000,largest_contentful_paint=2993400000,cumulative_layout_shift=0,first_input_delay=4200000,first_input_time=4488200000,time_spent=5953000000,in_foreground_periods="[{\"start\":2940200000,\"duration\":182400000},{\"start\":4491900000,\"duration\":1461100000}]",frustration_count=0 ___TIMESTAMP___
view,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=115.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4921,view_url=http://127.0.0.1:8081/index2.html,view_host=127.0.0.1:8081,view_path=/index2.html,view_name=/index2.html,view_path_group=/index2.html,view_url_query={},view_loading_type=initial_load is_active=true,view_error_count=1,view_resource_count=10,view_long_task_count=0,view_action_count=0,first_contentful_paint=2993400000,largest_contentful_paint=2993400000,cumulative_layout_shift=0,first_input_delay=4200000,first_input_time=4488200000,time_spent=5953000000,in_foreground_periods="[{\"start\":2940200000,\"duration\":182400000},{\"start\":4491900000,\"duration\":1461100000}]",frustration_count=0 ___TIMESTAMP___
view,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=116.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4922,view_url=http://127.0.0.1:8081/index3.html,view_host=127.0.0.1:8081,view_path=/index3.html,view_name=/index3.html,view_path_group=/index3.html,view_url_query={},view_loading_type=initial_load is_active=true,view_error_count=1,view_resource_count=10,view_long_task_count=0,view_action_count=0,first_contentful_paint=2993400000,largest_contentful_paint=2993400000,cumulative_layout_shift=0,first_input_delay=4200000,first_input_time=4488200000,time_spent=5953000000,in_foreground_periods="[{\"start\":2940200000,\"duration\":182400000},{\"start\":4491900000,\"duration\":1461100000}]",frustration_count=0 ___TIMESTAMP___
view,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=117.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4923,view_url=http://127.0.0.1:8081/index4.html,view_host=127.0.0.1:8081,view_path=/index4.html,view_name=/index4.html,view_path_group=/index4.html,view_url_query={},view_loading_type=initial_load is_active=true,view_error_count=1,view_resource_count=10,view_long_task_count=0,view_action_count=0,first_contentful_paint=2993400000,largest_contentful_paint=2993400000,cumulative_layout_shift=0,first_input_delay=4200000,first_input_time=4488200000,time_spent=5953000000,in_foreground_periods="[{\"start\":2940200000,\"duration\":182400000},{\"start\":4491900000,\"duration\":1461100000}]",frustration_count=0 ___TIMESTAMP___
`

var resourceBody = `resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://static.guance.com/browser-sdk/v3/dataflux-rum.js,resource_url_host=static.guance.com,resource_url_path=/browser-sdk/v3/dataflux-rum.js,resource_url_path_group=/browser-sdk/?/dataflux-rum.js,resource_url_query={},resource_type=js,resource_method=GET duration=2839600000,resource_size=147010,resource_tcp=495400000,resource_ssl=494900000,resource_ttfb=2084900000,resource_trans=251500000,resource_first_byte=2588100000,resource_download_time="{\"duration\":251500000,\"start\":2588100000}",resource_first_byte_time="{\"duration\":2084900000,\"start\":503200000}",resource_connect_time="{\"duration\":495400000,\"start\":7500000}",resource_ssl_time="{\"duration\":494900000,\"start\":8000000}" ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=http://127.0.0.1:8081/dist/bundle.js,resource_url_host=127.0.0.1,resource_url_path=/dist/bundle.js,resource_url_path_group=/dist/bundle.js,resource_url_query={},resource_type=js,resource_status=200,resource_status_group=2xx,resource_method=GET duration=8300000,resource_size=810,resource_ttfb=800000,resource_trans=100000,resource_first_byte=8200000,resource_download_time="{\"duration\":100000,\"start\":8200000}",resource_first_byte_time="{\"duration\":800000,\"start\":7400000}" ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://datadog-docs.imgix.net/images/real_user_monitoring/data_collected/event-hierarchy.9c1705008e5cb41df344aa7ed1374b3b.png,resource_url_host=datadog-docs.imgix.net,resource_url_path=/images/real_user_monitoring/data_collected/event-hierarchy.9c1705008e5cb41df344aa7ed1374b3b.png,resource_url_path_group=/images/real_user_monitoring/data_collected/?,resource_url_query={},resource_type=image,resource_method=GET duration=1508300000,resource_size=49304,resource_tcp=616800000,resource_ssl=616400000,resource_ttfb=649800000,resource_trans=226500000,resource_first_byte=1281800000,resource_download_time="{\"duration\":226500000,\"start\":1281800000}",resource_first_byte_time="{\"duration\":649800000,\"start\":632000000}",resource_connect_time="{\"duration\":616800000,\"start\":14900000}",resource_ssl_time="{\"duration\":616400000,\"start\":15300000}" ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=http://127.0.0.1:8081/index.html,resource_url_host=127.0.0.1,resource_url_path=/index.html,resource_url_path_group=/index.html,resource_url_query={},resource_type=document,resource_method=GET duration=62500000,resource_size=3377,resource_dns=0,resource_tcp=300000,resource_ttfb=56400000,resource_trans=200000,resource_redirect=0,resource_first_byte=56700000,resource_dns_time="{\"duration\":0,\"start\":5600000}",resource_download_time="{\"duration\":200000,\"start\":62300000}",resource_first_byte_time="{\"duration\":56400000,\"start\":5900000}",resource_connect_time="{\"duration\":300000,\"start\":5600000}",resource_redirect_time="{\"duration\":0,\"start\":0}" ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://cdn-dynmedia-1.microsoft.com/is/image/microsoftcorp/gldn-XSS-Hero-Xbox-Series-S:VP5-1596x600,resource_url_host=cdn-dynmedia-1.microsoft.com,resource_url_path=/is/image/microsoftcorp/gldn-XSS-Hero-Xbox-Series-S:VP5-1596x600,resource_url_path_group=/is/image/microsoftcorp/?,resource_url_query={},resource_type=image,resource_method=GET duration=4120000000 ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://www.fastly.com/cimages/6pk8mg3yh2ee/4zLuDTDqRq5H2Nuj8Ch0jB/d51cd0d20e0881043bec5c84151f1f4f/hero-laptop_1.png,resource_url_host=www.fastly.com,resource_url_path=/cimages/6pk8mg3yh2ee/4zLuDTDqRq5H2Nuj8Ch0jB/d51cd0d20e0881043bec5c84151f1f4f/hero-laptop_1.png,resource_url_path_group=/cimages/?/?/?/?,resource_url_query={},resource_type=image,resource_method=GET duration=4240100000 ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://is5-ssl.mzstatic.com/image/thumb/Music116/v4/b0/a5/b1/b0a5b100-e8a2-0a33-7835-4d8e1234d628/5054429157048.png/632x632bb.webp,resource_url_host=is5-ssl.mzstatic.com,resource_url_path=/image/thumb/Music116/v4/b0/a5/b1/b0a5b100-e8a2-0a33-7835-4d8e1234d628/5054429157048.png/632x632bb.webp,resource_url_path_group=/image/thumb/?/?/?/?/?/?/?/?,resource_url_query={},resource_type=image,resource_method=GET duration=4879200000,resource_size=27750,resource_tcp=1666500000,resource_ssl=1666100000,resource_ttfb=244400000,resource_trans=103300000,resource_first_byte=4775900000,resource_download_time="{\"duration\":103300000,\"start\":4775900000}",resource_first_byte_time="{\"duration\":244400000,\"start\":4531500000}",resource_connect_time="{\"duration\":1666500000,\"start\":2864900000}",resource_ssl_time="{\"duration\":1666100000,\"start\":2865300000}" ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://i.ytimg.com/vi/tXWfNvnKUug/hq720.jpg?sqp\=-oaymwEcCNAFEJQDSFXyq4qpAw4IARUAAIhCGAFwAcABBg\=\=&rs\=AOn4CLApZUDdPVnFtcTEorLSnZfmbnGmRg,resource_url_host=i.ytimg.com,resource_url_path=/vi/tXWfNvnKUug/hq720.jpg,resource_url_path_group=/vi/tXWfNvnKUug/?,resource_url_query={\"sqp\":\"-oaymwEcCNAFEJQDSFXyq4qpAw4IARUAAIhCGAFwAcABBg\=\=\"\,\"rs\":\"AOn4CLApZUDdPVnFtcTEorLSnZfmbnGmRg\"},resource_type=image,resource_method=GET duration=5344200000,resource_size=83720,resource_tcp=1544800000,resource_ssl=1544400000,resource_ttfb=683200000,resource_trans=251100000,resource_first_byte=5093100000,resource_download_time="{\"duration\":251100000,\"start\":5093100000}",resource_first_byte_time="{\"duration\":683200000,\"start\":4409900000}",resource_connect_time="{\"duration\":1544800000,\"start\":2864800000}",resource_ssl_time="{\"duration\":1544400000,\"start\":2865200000}" ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://pbs.twimg.com/media/Fey7tzdX0AIHpYk?format\=jpg&name\=medium,resource_url_host=pbs.twimg.com,resource_url_path=/media/Fey7tzdX0AIHpYk,resource_url_path_group=/media/?,resource_url_query={\"format\":\"jpg\"\,\"name\":\"medium\"},resource_type=image,resource_method=GET duration=5400900000 ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://www.cloudflare.com/static/cee7180ec61f2b210f14adc1b459063f/69740-Cloudflare-BDES-2164-Animations-Superside-01-D1_zh-CN.gif,resource_url_host=www.cloudflare.com,resource_url_path=/static/cee7180ec61f2b210f14adc1b459063f/69740-Cloudflare-BDES-2164-Animations-Superside-01-D1_zh-CN.gif,resource_url_path_group=/static/?/?,resource_url_query={},resource_type=image,resource_method=GET duration=5538800000 ___TIMESTAMP___
resource,sdk_name=df_web_rum_sdk,sdk_version=3.0.14,app_id=___APPID___,env=___ENV___,service=web-rum-demo,version=___VERSION___,userid=1e7e8237-3e69-4c71-99e9-078e0e825a62,session_id=df2ca679-5548-4ff2-b77d-a73d8f0d0a3c,session_type=user,is_signin=F,os=Mac\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=114.0.0.0,browser_version_major=114,screen_size=1440*900,network_type=wifi,view_id=cd88281d-5219-4856-8eb0-22aa05cd4920,view_url=http://127.0.0.1:8081/index.html,view_host=127.0.0.1:8081,view_path=/index.html,view_name=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://www.akamai.com/site/en/images/card/recording-basketball-game.jpg,resource_url_host=www.akamai.com,resource_url_path=/site/en/images/card/recording-basketball-game.jpg,resource_url_path_group=/site/en/images/card/recording-basketball-game.jpg,resource_url_query={},resource_type=image,resource_method=GET duration=8049000000 ___TIMESTAMP___
resource,sdk_name=df_android_rum_sdk,sdk_version=2.0.26,app_id=___APPID___,env=___ENV___,service=browser,version=___VERSION___,userid=6931fa8d-769b-46ef-998f-34fc20947562,session_id=49c6f9c2-318b-4c99-957d-b514db8706ee,session_type=user,is_signin=F,os=Mac\\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=103.0.0.0,browser_version_major=103,screen_size=1440*900,network_type=4g,view_id=8a5fcfb3-9f18-465e-8598-5adcc9924abc,view_url=http://localhost:8080/index.html?1111,view_host=localhost:8080,view_path=/index.html,view_path_group=/index.html,view_url_query={},resource_url=https://static.guance.com/browser-sdk/v2/dataflux-rum.js,resource_url_host=static.guance.com,resource_url_path=/browser-sdk/v2/dataflux-rum.js,resource_url_path_group=/browser-sdk/?/dataflux-rum.js,resource_url_query={},resource_type=js,resource_method=GET duration=0 ___TIMESTAMP___
resource,sdk_name=df_android_rum_sdk,sdk_version=2.0.26,app_id=___APPID___,env=___ENV___,service=browser,version=___VERSION___,userid=6931fa8d-769b-46ef-998f-34fc20947562,session_id=49c6f9c2-318b-4c99-957d-b514db8706ee,session_type=user,is_signin=F,os=Mac\\ OS,os_version=10.15.7,os_version_major=10,browser=Chrome,browser_version=103.0.0.0,browser_version_major=103,screen_size=1440*900,network_type=4g,view_id=8a5fcfb3-9f18-465e-8598-5adcc9924abc,view_url=http://localhost:8080/index.html?1111,view_host=localhost:8080,view_path=/index.html,view_path_group=/index.html,view_url_query={},resource_url=http://localhost:8080/dist/bundle.js,resource_url_host=localhost,resource_url_path=/dist/bundle.js,resource_url_path_group=/dist/bundle.js,resource_url_query={},resource_type=js,resource_method=GET duration=0 ___TIMESTAMP___`

var data = map[RUMType]string{
	RUMResource: resourceBody,
	RUMError:    errorBody,
	RUMView:     viewBody,
}

type counter struct {
	ok      int
	failure int
}

func doSend(ctx context.Context, typ RUMType) {

	cnt := &counter{}

	pts := data[typ]

	pts = strings.ReplaceAll(pts, "___APPID___", AppID)
	pts = strings.ReplaceAll(pts, "___ENV___", AppEnv)
	pts = strings.ReplaceAll(pts, "___VERSION___", AppVersion)
	pts = strings.ReplaceAll(pts, "___TIMESTAMP___", strconv.FormatInt(time.Now().UnixNano(), 10))

	fmt.Println(pts)

	log.Printf("body length: %d\n", len(pts))

	body := strings.NewReader(pts)

	for {

		select {
		case <-ctx.Done():
			outputChan <- cnt
			return
		default:
			if _, err := body.Seek(0, io.SeekStart); err != nil {
				log.Printf("seek error: %s", err)
				continue
			}
			resp, err := http.Post(rumEndpoint, "text/plain;charset=UTF-8", body)
			if err != nil {
				cnt.failure++
				log.Println(err)

			} else {
				func(resp *http.Response) {
					defer resp.Body.Close()

					if resp.StatusCode/100 == 2 {
						cnt.ok++
					} else {
						cnt.failure++
					}
				}(resp)
			}

		}

	}

}

func main() {

	ctx, f := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	for i := 0; i < concurrentCnt; i++ {

		go func(ctx context.Context) {

			<-flagChan
			wg.Add(1)
			doSend(ctx, RUMResource)
			wg.Done()

		}(ctx)
	}

	exitChan := make(chan struct{}, 1)

	go func() {

		sum := &counter{}

		for cnt := range outputChan {
			sum.ok += cnt.ok
			sum.failure += cnt.failure
		}

		log.Printf("ok count: %d    failure count: %d,   sum: %d in %d seconds\n", sum.ok, sum.failure, sum.ok+sum.failure, Duration/time.Second)
		close(exitChan)
	}()

	close(flagChan)
	time.Sleep(Duration)
	f()

	wg.Wait()
	close(outputChan)
	<-exitChan
}
