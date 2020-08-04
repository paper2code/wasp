package wasptest

import (
	"testing"
	"time"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	waspapi "github.com/iotaledger/wasp/packages/apilib"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/vm/examples/wasmpoc"
	"github.com/iotaledger/wasp/packages/vm/vmconst"
)

// sending 5 NOP requests with 1 sec sleep between each
func TestWasmVMSend5Requests1Sec(t *testing.T) {
	// setup
	wasps := setup(t, "test_cluster", "TestWasmVMSend5Requests1Sec")

	err := wasps.ListenToMessages(map[string]int{
		"bootuprec":           wasps.NumSmartContracts(),
		"active_committee":    1,
		"dismissed_committee": 0,
		"request_in":          6,
		"request_out":         7,
		"state":               -1, // must be 6 or 7
		"vmmsg":               -1,
	})
	check(err, t)

	_, err = PutBootupRecords(wasps)
	check(err, t)

	// number 5 is "Wasm VM PoC program" in cluster.json
	sc := &wasps.SmartContractConfig[4]

	err = Activate1SC(wasps, sc)
	check(err, t)

	err = CreateOrigin1SC(wasps, sc)
	check(err, t)

	scAddress := sc.SCAddress()
	ownerAddress := sc.OwnerAddress()

	for i := 0; i < 5; i++ {
		err = SendSimpleRequest(wasps, sc.OwnerSigScheme(), waspapi.CreateSimpleRequestParams{
			SCAddress:   &scAddress,
			RequestCode: wasmpoc.RequestNOP,
		})
		check(err, t)
		time.Sleep(1 * time.Second)
	}

	wasps.CollectMessages(10 * time.Second)

	if !wasps.Report() {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(scAddress, 1, map[balance.Color]int64{
		balance.ColorIOTA: 0,
		sc.GetColor():     1,
	}) {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(ownerAddress, testutil.RequestFundsAmount-1, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
	}) {
		t.Fail()
	}

	if !wasps.VerifySCState(sc, 0, map[kv.Key][]byte{
		vmconst.VarNameOwnerAddress: sc.GetColor().Bytes(),
		vmconst.VarNameProgramHash:  []byte(wasmpoc.ProgramHash),
	}) {
		t.Fail()
	}
}
