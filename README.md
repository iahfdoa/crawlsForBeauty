<h4 align="center">crawlsForBeauty 是一个用Go编写的美女图片爬取工具</h4>

```text
root@VM-4-13-ubuntu:~# crawlsForBeauty

                            __     ______           ____                   __
  ______________ __      __/ /____/ ____/___  _____/ __ )___  ____ ___  __/ /___  __
 / ___/ ___/ __ / | /| / / / ___/ /_  / __ \/ ___/ __  / _ \/ __ / / / / __/ / / /
/ /__/ /  / /_/ /| |/ |/ / (__  ) __/ / /_/ / /  / /_/ /  __/ /_/ / /_/ / /_/ /_/ /
\___/_/   \__,_/ |__/|__/_/____/_/    \____/_/  /_____/\___/\__,_/\__,_/\__/\__, /    v1.0.0
                                                                           /____/

慎用。你要为自己的行为负责
开发者不承担任何责任，也不对任何误用或损坏负责.

1s [======>-------------------------------------------------------------------------------------------]   8%
```

# 用法
```shell
crawlsForBeauty -h
```

```go
// tag 类型 t 为图库 i为tag
tagFunc := func(t, i int) string {
		switch t {
		case 2:
			switch i {
			case 1:
				return "latest"
			case 2:
				return "hot"
			case 3:
				return "toplist"
			case 4:
				return "random"
			default:
				return ""
			}
		case 1:
			switch i {
			case 1:
				return "stockings-porn"
			case 2:
				return "foot-fetish-porn"
			case 3:
				return "housewife-porn"
			case 4:
				return "teacher-porn"
			case 5:
				return "teen-porn"
			case 6:
				return "masturbation-porn"
			default:
				return "homemade-porn"
			}
		default:
			switch i {
			case 1:
				return "meitui"
			case 2:
				return "meixiong"
			case 3:
				return "meitun"
			case 4:
				return "shenyan"
			case 5:
				return "xiaoxinggan"
			case 6:
				return "xiaotianmei"
			default:
				return ""
			}

		}
	}
```
