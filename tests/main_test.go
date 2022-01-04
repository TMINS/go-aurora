package tests

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/awensir/go-aurora/aurora"
	"github.com/awensir/go-aurora/aurora/pprofs"
	"github.com/spf13/viper"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"os"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"time"
)

type Mm1 struct {
}

func (m *Mm1) PreHandle(c *aurora.Ctx) bool {
	log.Println("PreHandle Mm1")
	return true
}

func (m *Mm1) PostHandle(c *aurora.Ctx) {
	log.Println("PostHandle Mm1")
}

func (m *Mm1) AfterCompletion(c *aurora.Ctx) {
	log.Println("AfterCompletion Mm1")
}

type Mm2 struct {
}

func (m *Mm2) PreHandle(c *aurora.Ctx) bool {
	log.Println("PreHandle Mm2")

	return true
}

func (m *Mm2) PostHandle(c *aurora.Ctx) {
	log.Println("PostHandle Mm2")

}

func (m *Mm2) AfterCompletion(c *aurora.Ctx) {
	log.Println("AfterCompletion Mm2")
}

type Mm3 struct {
}

func (m *Mm3) PreHandle(c *aurora.Ctx) bool {
	log.Println("PreHandle Mm3")
	return true
}

func (m *Mm3) PostHandle(c *aurora.Ctx) {
	log.Println("PostHandle Mm3")
}

func (m *Mm3) AfterCompletion(c *aurora.Ctx) {
	log.Println("AfterCompletion Mm3")
}

// 拦截器测试
func TestIntercept(t *testing.T) {
	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		log.Println("service..")
		return nil
	})

	a.RouteIntercept("/", &Mm1{}, &Mm2{}, &Mm3{})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}

// TestWebSocketClient 发起socket 测试客户端
func TestWebSocketClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:8080/", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "内部错误！")

	err = wsjson.Write(ctx, c, "Hello WebSocket Server")
	if err != nil {
		panic(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("接收到服务端响应：%v\n", v)

	c.Close(websocket.StatusNormalClosure, "")
}

// TestWebSocketServer websocket 服务端测试
func TestWebSocketServer(t *testing.T) {
	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		accept, err := websocket.Accept(c.Response, c.Request, nil)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
		defer cancel()
		var v interface{}
		err = wsjson.Read(ctx, accept, &v)
		if err != nil {
			return err
		}
		log.Printf("接收到客户端：%v\n", v)
		err = wsjson.Write(ctx, accept, "Hello WebSocket Client")
		if err != nil {
			return err
		}
		accept.Close(websocket.StatusNormalClosure, "")
		return nil
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

// TestPprof 接口执行性能测试
func TestPprof(t *testing.T) {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		return nil
	})
	a.GET("/debug/pprof/heap", pprofs.Index)
	a.GET("/debug/pprof/cmdline", pprofs.Cmdline)
	a.GET("/debug/pprof/profile", pprofs.Profile)
	a.GET("/debug/pprof/symbol", pprofs.Symbol)
	a.GET("/debug/pprof/trace", pprofs.Trace)
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}

// TestPlugins 插件中断测试
func TestPlugins(t *testing.T) {
	//获取 aurora 路由实例
	a := aurora.New()
	a.Plugin(func(ctx *aurora.Ctx) bool {
		fmt.Println("1")
		return true
	}, func(ctx *aurora.Ctx) bool {
		fmt.Println("2")
		return true
	}, func(ctx *aurora.Ctx) bool {
		fmt.Println("3")
		return true
	}, func(ctx *aurora.Ctx) bool {
		fmt.Println("4")
		return true
	}, func(ctx *aurora.Ctx) bool {
		fmt.Println("5")
		ctx.Message("plugin error test!")
		return false
	})
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		return nil
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}

func TestColor(t *testing.T) {
	fmt.Println("")

	// 前景 背景 颜色
	// ---------------------------------------
	// 30  40  黑色
	// 31  41  红色
	// 32  42  绿色
	// 33  43  黄色
	// 34  44  蓝色
	// 35  45  紫红色
	// 36  46  青蓝色
	// 37  47  白色
	//
	// 代码 意义
	// -------------------------
	//  0  终端默认设置
	//  1  高亮显示
	//  4  使用下划线
	//  5  闪烁
	//  7  反白显示
	//  8  不可见

	for b := 40; b <= 47; b++ { // 背景色彩 = 40-47
		for f := 30; f <= 37; f++ { // 前景色彩 = 30-37
			for d := range []int{0, 1, 4, 5, 7, 8} { // 显示方式 = 0,1,4,5,7,8
				d += 1
				fmt.Printf(" %c[%d;%d;%dm%s(d=%d,b=%d,f=%d)%c[0m ", 0x1B, 1, 0, f, "", 1, 0, f, 0x1B)
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func TestLog(t *testing.T) {
	newLog := aurora.NewLog()

	go func() {
		for i := 0; i < 100; i++ {
			for i := 0; i < 100; i++ {
				newLog.Info("Info ", i)
				newLog.Warning("Warning ", i)
				newLog.Debug("Debug ", i)
				newLog.Error("Error ", i)
			}
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			newLog.Info("Info ", i)
			newLog.Warning("Warning ", i)
			newLog.Debug("Debug ", i)
			newLog.Error("Error ", i)
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			newLog.Info("Info ", i)
			newLog.Warning("Warning ", i)
			newLog.Debug("Debug ", i)
			newLog.Error("Error ", i)
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			newLog.Info("Info ", i)
			newLog.Warning("Warning ", i)
			newLog.Debug("Debug ", i)
			newLog.Error("Error ", i)
		}
	}()

	time.Sleep(5 * time.Second)
}

func TestJ(t *testing.T) {
	//s := "{\n    \"name\":\"test\",\n    \"age\":15,\n    \"gender\":\"nv\"\n}"
	s := "23 a56"
	compile := regexp.MustCompile("\\d*")
	find := compile.FindAllString(s, -1)
	fmt.Println(find)
}

func TestConfigFile(t *testing.T) {
	v := viper.New()
	v.SetConfigType("json")
	err := v.ReadConfig(bytes.NewBuffer(static))
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	static := v.Get("type")
	fmt.Println(static)
}

func TestPath(t *testing.T) {
	//url := "/sada/{sss}/{bbb}"

}
func Check(url string) {
	if strings.Count(url, "{") == strings.Count(url, "}") {

	}
}

var static = []byte(`{
  "type": {
    ".323": "text/h323",
    ".3gp": "video/3gpp",
    ".aab": "application/x-authoware-bin",
    ".aam": "application/x-authoware-map",
    ".aas": "application/x-authoware-seg",
    ".acx": "application/internet-property-stream",
    ".ai": "application/postscript",
    ".aif": "audio/x-aiff",
    ".aifc": "audio/x-aiff",
    ".aiff": "audio/x-aiff",
    ".als": "audio/X-Alpha5",
    ".amc": "application/x-mpeg",
    ".ani": "application/octet-stream",
    ".apk": "application/vnd.android.package-archive",
    ".asc": "text/plain",
    ".asd": "application/astound",
    ".asf": "video/x-ms-asf",
    ".asn": "application/astound",
    ".asp": "application/x-asap",
    ".asr": "video/x-ms-asf",
    ".asx": "video/x-ms-asf",
    ".au": "audio/basic",
    ".avb": "application/octet-stream",
    ".avi": "video/x-msvideo",
    ".awb": "audio/amr-wb",
    ".axs": "application/olescript",
    ".bas": "text/plain",
    ".bcpio": "application/x-bcpio",
    ".bin ": "application/octet-stream",
    ".bld": "application/bld",
    ".bld2": "application/bld2",
    ".bmp": "image/bmp",
    ".bpk": "application/octet-stream",
    ".bz2": "application/x-bzip2",
    ".c": "text/plain",
    ".cal": "image/x-cals",
    ".cat": "application/vnd.ms-pkiseccat",
    ".ccn": "application/x-cnc",
    ".cco": "application/x-cocoa",
    ".cdf": "application/x-cdf",
    ".cer": "application/x-x509-ca-cert",
    ".cgi": "magnus-internal/cgi",
    ".chat": "application/x-chat",
    ".class": "application/octet-stream",
    ".clp": "application/x-msclip",
    ".cmx": "image/x-cmx",
    ".co": "application/x-cult3d-object",
    ".cod": "image/cis-cod",
    ".conf": "text/plain",
    ".cpio": "application/x-cpio",
    ".cpp": "text/plain",
    ".cpt": "application/mac-compactpro",
    ".crd": "application/x-mscardfile",
    ".crl": "application/pkix-crl",
    ".crt": "application/x-x509-ca-cert",
    ".csh": "application/x-csh",
    ".csm": "chemical/x-csml",
    ".csml": "chemical/x-csml",
    ".css": "text/css",
    ".cur": "application/octet-stream",
    ".dcm": "x-lml/x-evm",
    ".dcr": "application/x-director",
    ".dcx": "image/x-dcx",
    ".der": "application/x-x509-ca-cert",
    ".dhtml": "text/html",
    ".dir": "application/x-director",
    ".dll": "application/x-msdownload",
    ".dmg": "application/octet-stream",
    ".dms": "application/octet-stream",
    ".doc": "application/msword",
    ".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    ".dot": "application/msword",
    ".dvi": "application/x-dvi",
    ".dwf": "drawing/x-dwf",
    ".dwg": "application/x-autocad",
    ".dxf": "application/x-autocad",
    ".dxr": "application/x-director",
    ".ebk": "application/x-expandedbook",
    ".emb": "chemical/x-embl-dl-nucleotide",
    ".embl": "chemical/x-embl-dl-nucleotide",
    ".eps": "application/postscript",
    ".epub": "application/epub+zip",
    ".eri": "image/x-eri",
    ".es": "audio/echospeech",
    ".esl": "audio/echospeech",
    ".etc": "application/x-earthtime",
    ".etx": "text/x-setext",
    ".evm": "x-lml/x-evm",
    ".evy": "application/envoy",
    ".exe": "application/octet-stream",
    ".fh4": "image/x-freehand",
    ".fh5": "image/x-freehand",
    ".fhc": "image/x-freehand",
    ".fif": "application/fractals",
    ".flr": "x-world/x-vrml",
    ".flv": "flv-application/octet-stream",
    ".fm": "application/x-maker",
    ".fpx": "image/x-fpx",
    ".fvi": "video/isivideo",
    ".gau": "chemical/x-gaussian-input",
    ".gca": "application/x-gca-compressed",
    ".gdb": "x-lml/x-gdb",
    ".gif": "image/gif",
    ".gps": "application/x-gps",
    ".gtar": "application/x-gtar",
    ".gz": "application/x-gzip",
    ".h": "text/plain",
    ".hdf": "application/x-hdf",
    ".hdm": "text/x-hdml",
    ".hdml": "text/x-hdml",
    ".hlp": "application/winhlp",
    ".hqx": "application/mac-binhex40",
    ".hta": "application/hta",
    ".htc": "text/x-component",
    ".htm": "text/html",
    ".html": "text/html",
    ".hts": "text/html",
    ".htt": "text/webviewhtml",
    ".ice": "x-conference/x-cooltalk",
    ".ico": "image/x-icon",
    ".ief": "image/ief",
    ".ifm": "image/gif",
    ".ifs": "image/ifs",
    ".iii": "application/x-iphone",
    ".imy": "audio/melody",
    ".ins": "application/x-internet-signup",
    ".ips": "application/x-ipscript",
    ".ipx": "application/x-ipix",
    ".isp": "application/x-internet-signup",
    ".it": "audio/x-mod",
    ".itz": "audio/x-mod",
    ".ivr": "i-world/i-vrml",
    ".j2k": "image/j2k",
    ".jad": "text/vnd.sun.j2me.app-descriptor",
    ".jam": "application/x-jam",
    ".jar": "application/java-archive",
    ".java": "text/plain",
    ".jfif": "image/pipeg",
    ".jnlp": "application/x-java-jnlp-file",
    ".jpe": "image/jpeg",
    ".jpeg": "image/jpeg",
    ".jpg": "image/jpeg",
    ".jpz": "image/jpeg",
    ".js": "application/x-javascript",
    ".jwc": "application/jwc",
    ".kjx": "application/x-kjx",
    ".lak": "x-lml/x-lak",
    ".latex": "application/x-latex",
    ".lcc": "application/fastman",
    ".lcl": "application/x-digitalloca",
    ".lcr": "application/x-digitalloca",
    ".lgh": "application/lgh",
    ".lha": "application/octet-stream",
    ".lml": "x-lml/x-lml",
    ".lmlpack": "x-lml/x-lmlpack",
    ".log": "text/plain",
    ".lsf": "video/x-la-asf",
    ".lsx": "video/x-la-asf",
    ".lzh": "application/octet-stream",
    ".m13": "application/x-msmediaview",
    ".m14": "application/x-msmediaview",
    ".m15": "audio/x-mod",
    ".m3u": "audio/x-mpegurl",
    ".m3url": "audio/x-mpegurl",
    ".m4a": "audio/mp4a-latm",
    ".m4b": "audio/mp4a-latm",
    ".m4p": "audio/mp4a-latm",
    ".m4u": "video/vnd.mpegurl",
    ".m4v": "video/x-m4v",
    ".ma1": "audio/ma1",
    ".ma2": "audio/ma2",
    ".ma3": "audio/ma3",
    ".ma5": "audio/ma5",
    ".man": "application/x-troff-man",
    ".map": "magnus-internal/imagemap",
    ".mbd": "application/mbedlet",
    ".mct": "application/x-mascot",
    ".mdb": "application/x-msaccess",
    ".mdz": "audio/x-mod",
    ".me": "application/x-troff-me",
    ".mel": "text/x-vmel",
    ".mht": "message/rfc822",
    ".mhtml": "message/rfc822",
    ".mi": "application/x-mif",
    ".mid": "audio/mid",
    ".midi": "audio/midi",
    ".mif": "application/x-mif",
    ".mil": "image/x-cals",
    ".mio": "audio/x-mio",
    ".mmf": "application/x-skt-lbs",
    ".mng": "video/x-mng",
    ".mny": "application/x-msmoney",
    ".moc": "application/x-mocha",
    ".mocha": "application/x-mocha",
    ".mod": "audio/x-mod",
    ".mof": "application/x-yumekara",
    ".mol": "chemical/x-mdl-molfile",
    ".mop": "chemical/x-mopac-input",
    ".mov": "video/quicktime",
    ".movie": "video/x-sgi-movie",
    ".mp2": "video/mpeg",
    ".mp3": "audio/mpeg",
    ".mp4": "video/mp4",
    ".mpa": "video/mpeg",
    ".mpc": "application/vnd.mpohun.certificate",
    ".mpe": "video/mpeg",
    ".mpeg": "video/mpeg",
    ".mpg": "video/mpeg",
    ".mpg4": "video/mp4",
    ".mpga": "audio/mpeg",
    ".mpn": "application/vnd.mophun.application",
    ".mpp": "application/vnd.ms-project",
    ".mps": "application/x-mapserver",
    ".mpv2": "video/mpeg",
    ".mrl": "text/x-mrml",
    ".mrm": "application/x-mrm",
    ".ms": "application/x-troff-ms",
    ".msg": "application/vnd.ms-outlook",
    ".mts": "application/metastream",
    ".mtx": "application/metastream",
    ".mtz": "application/metastream",
    ".mvb": "application/x-msmediaview",
    ".mzv": "application/metastream",
    ".nar": "application/zip",
    ".nbmp": "image/nbmp",
    ".nc": "application/x-netcdf",
    ".ndb": "x-lml/x-ndb",
    ".ndwn": "application/ndwn",
    ".nif": "application/x-nif",
    ".nmz": "application/x-scream",
    ".nokia-op-logo": "image/vnd.nok-oplogo-color",
    ".npx": "application/x-netfpx",
    ".nsnd": "audio/nsnd",
    ".nva": "application/x-neva1",
    ".nws": "message/rfc822",
    ".oda": "application/oda",
    ".ogg": "audio/ogg",
    ".oom": "application/x-AtlasMate-Plugin",
    ".p10": "application/pkcs10",
    ".p12": "application/x-pkcs12",
    ".p7b": "application/x-pkcs7-certificates",
    ".p7c": "application/x-pkcs7-mime",
    ".p7m": "application/x-pkcs7-mime",
    ".p7r": "application/x-pkcs7-certreqresp",
    ".p7s": "application/x-pkcs7-signature",
    ".pac": "audio/x-pac",
    ".pae": "audio/x-epac",
    ".pan": "application/x-pan",
    ".pbm": "image/x-portable-bitmap",
    ".pcx": "image/x-pcx",
    ".pda": "image/x-pda",
    ".pdb": "chemical/x-pdb",
    ".pdf": "application/pdf",
    ".pfr": "application/font-tdpfr",
    ".pfx": "application/x-pkcs12",
    ".pgm": "image/x-portable-graymap",
    ".pict": "image/x-pict",
    ".pko": "application/ynd.ms-pkipko",
    ".pm": "application/x-perl",
    ".pma": "application/x-perfmon",
    ".pmc": "application/x-perfmon",
    ".pmd": "application/x-pmd",
    ".pml": "application/x-perfmon",
    ".pmr": "application/x-perfmon",
    ".pmw": "application/x-perfmon",
    ".png": "image/png",
    ".pnm": "image/x-portable-anymap",
    ".pnz": "image/png",
    ".pot,": "application/vnd.ms-powerpoint",
    ".ppm": "image/x-portable-pixmap",
    ".pps": "application/vnd.ms-powerpoint",
    ".ppt": "application/vnd.ms-powerpoint",
    ".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    ".pqf": "application/x-cprplayer",
    ".pqi": "application/cprplayer",
    ".prc": "application/x-prc",
    ".prf": "application/pics-rules",
    ".prop": "text/plain",
    ".proxy": "application/x-ns-proxy-autoconfig",
    ".ps": "application/postscript",
    ".ptlk": "application/listenup",
    ".pub": "application/x-mspublisher",
    ".pvx": "video/x-pv-pvx",
    ".qcp": "audio/vnd.qcelp",
    ".qt": "video/quicktime",
    ".qti": "image/x-quicktime",
    ".qtif": "image/x-quicktime",
    ".r3t": "text/vnd.rn-realtext3d",
    ".ra": "audio/x-pn-realaudio",
    ".ram": "audio/x-pn-realaudio",
    ".rar": "application/octet-stream",
    ".ras": "image/x-cmu-raster",
    ".rc": "text/plain",
    ".rdf": "application/rdf+xml",
    ".rf": "image/vnd.rn-realflash",
    ".rgb": "image/x-rgb",
    ".rlf": "application/x-richlink",
    ".rm": "audio/x-pn-realaudio",
    ".rmf": "audio/x-rmf",
    ".rmi": "audio/mid",
    ".rmm": "audio/x-pn-realaudio",
    ".rmvb": "audio/x-pn-realaudio",
    ".rnx": "application/vnd.rn-realplayer",
    ".roff": "application/x-troff",
    ".rp": "image/vnd.rn-realpix",
    ".rpm": "audio/x-pn-realaudio-plugin",
    ".rt": "text/vnd.rn-realtext",
    ".rte": "x-lml/x-gps",
    ".rtf": "application/rtf",
    ".rtg": "application/metastream",
    ".rtx": "text/richtext",
    ".rv": "video/vnd.rn-realvideo",
    ".rwc": "application/x-rogerwilco",
    ".s3m": "audio/x-mod",
    ".s3z": "audio/x-mod",
    ".sca": "application/x-supercard",
    ".scd": "application/x-msschedule",
    ".sct": "text/scriptlet",
    ".sdf": "application/e-score",
    ".sea": "application/x-stuffit",
    ".setpay": "application/set-payment-initiation",
    ".setreg": "application/set-registration-initiation",
    ".sgm": "text/x-sgml",
    ".sgml": "text/x-sgml",
    ".sh": "application/x-sh",
    ".shar": "application/x-shar",
    ".shtml": "magnus-internal/parsed-html",
    ".shw": "application/presentations",
    ".si6": "image/si6",
    ".si7": "image/vnd.stiwap.sis",
    ".si9": "image/vnd.lgtwap.sis",
    ".sis": "application/vnd.symbian.install",
    ".sit": "application/x-stuffit",
    ".skd": "application/x-Koan",
    ".skm": "application/x-Koan",
    ".skp": "application/x-Koan",
    ".skt": "application/x-Koan",
    ".slc": "application/x-salsa",
    ".smd": "audio/x-smd",
    ".smi": "application/smil",
    ".smil": "application/smil",
    ".smp": "application/studiom",
    ".smz": "audio/x-smd",
    ".snd": "audio/basic",
    ".spc": "application/x-pkcs7-certificates",
    ".spl": "application/futuresplash",
    ".spr": "application/x-sprite",
    ".sprite": "application/x-sprite",
    ".sdp": "application/sdp",
    ".spt": "application/x-spt",
    ".src": "application/x-wais-source",
    ".sst": "application/vnd.ms-pkicertstore",
    ".stk": "application/hyperstudio",
    ".stl": "application/vnd.ms-pkistl",
    ".stm": "text/html",
    ".svg": "image/svg+xml",
    ".sv4cpio": "application/x-sv4cpio",
    ".sv4crc": "application/x-sv4crc",
    ".svf": "image/vnd",
    ".svg": "image/svg+xml",
    ".svh": "image/svh",
    ".svr": "x-world/x-svr",
    ".swf": "application/x-shockwave-flash",
    ".swfl": "application/x-shockwave-flash",
    ".t": "application/x-troff",
    ".tad": "application/octet-stream",
    ".talk": "text/x-speech",
    ".tar": "application/x-tar",
    ".taz": "application/x-tar",
    ".tbp": "application/x-timbuktu",
    ".tbt": "application/x-timbuktu",
    ".tcl": "application/x-tcl",
    ".tex": "application/x-tex",
    ".texi": "application/x-texinfo",
    ".texinfo": "application/x-texinfo",
    ".tgz": "application/x-compressed",
    ".thm": "application/vnd.eri.thm",
    ".tif": "image/tiff",
    ".tiff": "image/tiff",
    ".tki": "application/x-tkined",
    ".tkined": "application/x-tkined",
    ".toc": "application/toc",
    ".toy": "image/toy",
    ".tr": "application/x-troff",
    ".trk": "x-lml/x-gps",
    ".trm": "application/x-msterminal",
    ".tsi": "audio/tsplayer",
    ".tsp": "application/dsptype",
    ".tsv": "text/tab-separated-values",
    ".ttf": "application/octet-stream",
    ".ttz": "application/t-time",
    ".txt": "text/plain",
    ".uls": "text/iuls",
    ".ult": "audio/x-mod",
    ".ustar": "application/x-ustar",
    ".uu": "application/x-uuencode",
    ".uue": "application/x-uuencode",
    ".vcd": "application/x-cdlink",
    ".vcf": "text/x-vcard",
    ".vdo": "video/vdo",
    ".vib": "audio/vib",
    ".viv": "video/vivo",
    ".vivo": "video/vivo",
    ".vmd": "application/vocaltec-media-desc",
    ".vmf": "application/vocaltec-media-file",
    ".vmi": "application/x-dreamcast-vms-info",
    ".vms": "application/x-dreamcast-vms",
    ".vox": "audio/voxware",
    ".vqe": "audio/x-twinvq-plugin",
    ".vqf": "audio/x-twinvq",
    ".vql": "audio/x-twinvq",
    ".vre": "x-world/x-vream",
    ".vrml": "x-world/x-vrml",
    ".vrt": "x-world/x-vrt",
    ".vrw": "x-world/x-vream",
    ".vts": "workbook/formulaone",
    ".wav": "audio/x-wav",
    ".wax": "audio/x-ms-wax",
    ".wbmp": "image/vnd.wap.wbmp",
    ".wcm": "application/vnd.ms-works",
    ".wdb": "application/vnd.ms-works",
    ".web": "application/vnd.xara",
    ".wi": "image/wavelet",
    ".wis": "application/x-InstallShield",
    ".wks": "application/vnd.ms-works",
    ".wm": "video/x-ms-wm",
    ".wma": "audio/x-ms-wma",
    ".wmd": "application/x-ms-wmd",
    ".wmf": "application/x-msmetafile",
    ".wml": "text/vnd.wap.wml",
    ".wmlc": "application/vnd.wap.wmlc",
    ".wmls": "text/vnd.wap.wmlscript",
    ".wmlsc": "application/vnd.wap.wmlscriptc",
    ".wmlscript": "text/vnd.wap.wmlscript",
    ".wmv": "audio/x-ms-wmv",
    ".wmx": "video/x-ms-wmx",
    ".wmz": "application/x-ms-wmz",
    ".wpng": "image/x-up-wpng",
    ".wps": "application/vnd.ms-works",
    ".wpt": "x-lml/x-gps",
    ".wri": "application/x-mswrite",
    ".wrl": "x-world/x-vrml",
    ".wrz": "x-world/x-vrml",
    ".ws": "text/vnd.wap.wmlscript",
    ".wsc": "application/vnd.wap.wmlscriptc",
    ".wv": "video/wavelet",
    ".wvx": "video/x-ms-wvx",
    ".wxl": "application/x-wxl",
    ".x-gzip": "application/x-gzip",
    ".xaf": "x-world/x-vrml",
    ".xar": "application/vnd.xara",
    ".xbm": "image/x-xbitmap",
    ".xdm": "application/x-xdma",
    ".xdma": "application/x-xdma",
    ".xdw": "application/vnd.fujixerox.docuworks",
    ".xht": "application/xhtml+xml",
    ".xhtm": "application/xhtml+xml",
    ".xhtml": "application/xhtml+xml",
    ".xla": "application/vnd.ms-excel",
    ".xlc": "application/vnd.ms-excel",
    ".xll": "application/x-excel",
    ".xlm": "application/vnd.ms-excel",
    ".xls": "application/vnd.ms-excel",
    ".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    ".xlt": "application/vnd.ms-excel",
    ".xlw": "application/vnd.ms-excel",
    ".xm": "audio/x-mod",
    ".xml": "text/plain",
    ".xml": "application/xml",
    ".xmz": "audio/x-mod",
    ".xof": "x-world/x-vrml",
    ".xpi": "application/x-xpinstall",
    ".xpm": "image/x-xpixmap",
    ".xsit": "text/xml",
    ".xsl": "text/xml",
    ".xul": "text/xul",
    ".xwd": "image/x-xwindowdump",
    ".xyz": "chemical/x-pdb",
    ".yz1": "application/x-yz1",
    ".z": "application/x-compress",
    ".zac": "application/x-zaurus-zac",
    ".zip": "application/zip",
    ".json": "application/json"
  }
}`)
