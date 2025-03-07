/*******************************************************************************
 * Copyright 2019-2020 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

// Package main provides C interface for orchestration
package main

/*
#include <stdlib.h>

#ifndef __ORCHESTRATION_H__
#define __ORCHESTRATION_H__

#ifdef __cplusplus
extern "C"
{
#endif

#define MAX_SVC_INFO_NUM 3
typedef struct {
	char* ExecutionType;
	char* ExeCmd;
} RequestServiceInfo;

typedef struct {
	char* ExecutionType;
	char* Target;
} TargetInfo;

typedef struct {
	char*      Message;
	char*      ServiceName;
	TargetInfo RemoteTargetInfo;
} ResponseService;

typedef char* (*identityGetterFunc)();
typedef char* (*keyGetterFunc)(char* id);

identityGetterFunc iGetter;
keyGetterFunc kGetter;

static void setHandler(identityGetterFunc ihandle, keyGetterFunc khandle){
	iGetter = ihandle;
	kGetter = khandle;
}

static char* bridge_iGetter(){
	return iGetter();
}

static char* bridge_kGetter(char* id){
	return kGetter(id);
}
#ifdef __cplusplus
}

#endif

#endif // __ORCHESTRATION_H__

*/
import "C"
import (
	"flag"
	"math"
	"strings"
	"sync"
	"unsafe"

	"github.com/lf-edge/edge-home-orchestration-go/internal/common/fscreator"
	"github.com/lf-edge/edge-home-orchestration-go/internal/common/logmgr"
	"github.com/lf-edge/edge-home-orchestration-go/internal/controller/configuremgr"
	"github.com/lf-edge/edge-home-orchestration-go/internal/controller/discoverymgr"
	mnedcmgr "github.com/lf-edge/edge-home-orchestration-go/internal/controller/discoverymgr/mnedc"
	scoringmgr "github.com/lf-edge/edge-home-orchestration-go/internal/controller/scoringmgr"
	"github.com/lf-edge/edge-home-orchestration-go/internal/controller/securemgr"
	"github.com/lf-edge/edge-home-orchestration-go/internal/controller/securemgr/verifier"
	"github.com/lf-edge/edge-home-orchestration-go/internal/controller/servicemgr"
	"github.com/lf-edge/edge-home-orchestration-go/internal/controller/servicemgr/executor/nativeexecutor"
	"github.com/lf-edge/edge-home-orchestration-go/internal/controller/storagemgr"
	"github.com/lf-edge/edge-home-orchestration-go/internal/db/bolt/wrapper"
	"github.com/lf-edge/edge-home-orchestration-go/internal/orchestrationapi"
	"github.com/lf-edge/edge-home-orchestration-go/internal/restinterface/cipher/dummy"
	"github.com/lf-edge/edge-home-orchestration-go/internal/restinterface/cipher/sha256"
	"github.com/lf-edge/edge-home-orchestration-go/internal/restinterface/client/restclient"
	"github.com/lf-edge/edge-home-orchestration-go/internal/restinterface/externalhandler"
	"github.com/lf-edge/edge-home-orchestration-go/internal/restinterface/internalhandler"
	"github.com/lf-edge/edge-home-orchestration-go/internal/restinterface/route"
	"github.com/lf-edge/edge-home-orchestration-go/internal/restinterface/tls"
)

const logPrefix = "[interface]"

// Handle Platform Dependencies
const (
	platform      = "linux"
	executionType = "native"

	edgeDir = "/var/edge-orchestration"

	logPath             = edgeDir + "/log"
	configPath          = edgeDir + "/apps"
	dbPath              = edgeDir + "/data/db"
	certificateFilePath = edgeDir + "/data/cert"

	cipherKeyFilePath = edgeDir + "/user/orchestration_userID.txt"
	deviceIDFilePath  = edgeDir + "/device/orchestration_deviceID.txt"
	mnedcServerConfig = edgeDir + "/mnedc/client-config.yaml"
)

var (
	flagVersion       bool
	commitID, version string
	log               = logmgr.GetInstance()

	orcheEngine orchestrationapi.Orche
)

// OrchestrationInit runs orchestration service and discovers remote orchestration services
//export OrchestrationInit
func OrchestrationInit(secure C.int, mnedc C.int) C.int {
	flag.BoolVar(&flagVersion, "v", false, "if true, print version and exit")
	flag.BoolVar(&flagVersion, "version", false, "if true, print version and exit")
	flag.Parse()

	logmgr.InitLogfile(logPath)
	log.Println(logPrefix, "OrchestrationInit")
	log.Println(">>> commitID  : ", commitID)
	log.Println(">>> version   : ", version)
	wrapper.SetBoltDBPath(dbPath)

	if err := fscreator.CreateFileSystem(edgeDir); err != nil {
		log.Panicf("%s Failed to create edge-orchestration file system\n", logPrefix)
		return -1
	}

	isSecured := false
	if secure == 1 {
		log.Println(logPrefix, "Orchestration init with secure option")
		securemgr.Start(edgeDir)
		isSecured = true
	}

	cipher := dummy.GetCipher(cipherKeyFilePath)
	if isSecured {
		cipher = sha256.GetCipher(cipherKeyFilePath)
	}

	restIns := restclient.GetRestClient()
	restIns.SetCipher(cipher)

	servicemgr.GetInstance().SetClient(restIns)
	discoverymgr.GetInstance().SetClient(restIns)

	builder := orchestrationapi.OrchestrationBuilder{}
	builder.SetWatcher(configuremgr.GetInstance(configPath, executionType))
	builder.SetDiscovery(discoverymgr.GetInstance())
	builder.SetStorage(storagemgr.GetInstance())
	builder.SetVerifierConf(verifier.GetInstance())
	builder.SetScoring(scoringmgr.GetInstance())
	builder.SetService(servicemgr.GetInstance())
	builder.SetExecutor(nativeexecutor.GetInstance())
	builder.SetClient(restIns)
	orcheEngine = builder.Build()
	if orcheEngine == nil {
		log.Fatalf("%s Orchestaration initialize fail", logPrefix)
		return -1
	}

	orcheEngine.Start(deviceIDFilePath, platform, executionType)

	var restEdgeRouter *route.RestRouter
	if isSecured {
		restEdgeRouter = route.NewRestRouterWithCerti(certificateFilePath)
	} else {
		restEdgeRouter = route.NewRestRouter()
	}

	internalapi, err := orchestrationapi.GetInternalAPI()
	if err != nil {
		log.Fatalf("%s Orchestaration internal api : %s", logPrefix, err.Error())
	}
	ihandle := internalhandler.GetHandler()
	ihandle.SetOrchestrationAPI(internalapi)

	if isSecured {
		ihandle.SetCertificateFilePath(certificateFilePath)
	}
	ihandle.SetCipher(cipher)
	restEdgeRouter.Add(ihandle)

	externalapi, err := orchestrationapi.GetExternalAPI()
	if err != nil {
		log.Fatalf("%s Orchestaration external api : %s", logPrefix, err.Error())
	}
	ehandle := externalhandler.GetHandler()
	ehandle.SetOrchestrationAPI(externalapi)
	ehandle.SetCipher(dummy.GetCipher(cipherKeyFilePath))
	restEdgeRouter.Add(ehandle)

	restEdgeRouter.Start()

	log.Println(logPrefix, "Orchestration init done")
	mnedcmgr.GetServerInstance().SetCipher(cipher)
	if isSecured {
		mnedcmgr.GetServerInstance().SetCertificateFilePath(certificateFilePath)
		mnedcmgr.GetClientInstance().SetCertificateFilePath(certificateFilePath)
	}
	isMNEDCServer := false
	isMNEDCClient := false
	if mnedc == 1 {
		isMNEDCServer = true
		log.Println(logPrefix, "Orchestration init with MNEDC server option")
	} else if mnedc == 2 {
		isMNEDCClient = true
		log.Println(logPrefix, "Orchestration init with MNEDC client option")
	}

	go func() {
		if isMNEDCServer {
			mnedcmgr.GetServerInstance().StartMNEDCServer(deviceIDFilePath)
		} else if isMNEDCClient {
			mnedcmgr.GetClientInstance().StartMNEDCClient(deviceIDFilePath, mnedcServerConfig)
		}
	}()
	return 0
}

// OrchestrationRequestService performs request from service applications which uses orchestration service
//export OrchestrationRequestService
func OrchestrationRequestService(cAppName *C.char, cSelfSelection C.int, cRequester *C.char, serviceInfo *C.RequestServiceInfo, count C.int) C.ResponseService {
	log.Printf("%s OrchestrationRequestService", logPrefix)

	appName := C.GoString(cAppName)

	requestInfos := make([]orchestrationapi.RequestServiceInfo, count)
	CServiceInfo := (*[(math.MaxInt16 - 1) / unsafe.Sizeof(serviceInfo)]C.RequestServiceInfo)(unsafe.Pointer(serviceInfo))[:count:count]

	for idx, requestInfo := range CServiceInfo {
		requestInfos[idx].ExecutionType = C.GoString(requestInfo.ExecutionType)

		args := strings.Split(C.GoString(requestInfo.ExeCmd), " ")
		if strings.Compare(args[0], "") == 0 {
			args = nil
		}
		requestInfos[idx].ExeCmd = append([]string{}, args...)
	}

	externalAPI, err := orchestrationapi.GetExternalAPI()
	if err != nil {
		log.Fatalf("%s Orchestaration external api : %s", logPrefix, err.Error())
	}

	selfSel := true
	if cSelfSelection == 0 {
		selfSel = false
	}

	requester := C.GoString(cRequester)

	log.Printf("[OrchestrationRequestService] appName:%s", appName)
	log.Printf("[OrchestrationRequestService] selfSel:%v", selfSel)
	log.Printf("[OrchestrationRequestService] requester:%s", requester)
	log.Printf("[OrchestrationRequestService] infos:%v", requestInfos)

	res := externalAPI.RequestService(orchestrationapi.ReqeustService{
		ServiceName:      appName,
		SelfSelection:    selfSel,
		ServiceInfo:      requestInfos,
		ServiceRequester: requester,
	})
	log.Println(logPrefix, "requestService handle : ", res)

	ret := C.ResponseService{}
	ret.Message = C.CString(res.Message)
	ret.ServiceName = C.CString(res.ServiceName)
	ret.RemoteTargetInfo.ExecutionType = C.CString(res.RemoteTargetInfo.ExecutionType)
	ret.RemoteTargetInfo.Target = C.CString(res.RemoteTargetInfo.Target)

	return ret
}

type customHandler struct{}

// SetHandler sets handler
//export SetHandler
func SetHandler(iGetter C.identityGetterFunc, kGetter C.keyGetterFunc) {
	C.setHandler(iGetter, kGetter)
	tls.SetHandler(customHandler{})
}

var count int
var mtx sync.Mutex

// PrintLog provides logging interface
//export PrintLog
func PrintLog(cMsg *C.char) (count C.int) {
	mtx.Lock()
	msg := C.GoString(cMsg)
	defer mtx.Unlock()
	log.Printf(msg)
	count++
	return
}

func main() {
	// Do nothing because it is only to build static lWatcher interface
}
