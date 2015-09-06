package v8dispatcher

//func TestConsole(t *testing.T) {
//	worker, dist := newWorkerAndDispatcher(t)
//	dist.Register("console", Console{})

//	src, err := ioutil.ReadFile("setup.js")

//	err = worker.Load("console.js", `
//		console.log("size",42);

//		function putit_togo(arg) {
//			go_dispatch(function_registry.none, "echo", "noreturn", arg);
//		}

//		function getit_fromgo(then) {
//			go_dispatch(function_registry.put(then), "echo", "return", 42);
//		}
//		getit_fromgo(function(msg){
//			$print(msg);
//		});

//		putit_togo(36)

//		go_dispatch(function_registry.none, "badthings", "happen", "today");
//	`)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
