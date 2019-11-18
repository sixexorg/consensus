package plots

import (
	"encoding/binary"
	"github.com/sixexorg/consensus/crypto/shabal"
	"math/big"
)

type PlotCalculator interface {
	CalculateGenerationSignature(lastGenSig []byte, lastGenId int64) []byte

	CalculateScoop(genSig []byte, height int64) int

	CalculateHit1(accountId int64, nonce int, genSig []byte, scoop, pocVersion int) *big.Int

	CalculateHit2(genSig, scoopData []byte) *big.Int

	CalculateDeadline(accountId int64, nonce int, genSig []byte, scoop int, baseTarget int64, pocVersion int) *big.Int
}

type PlotCalculatorImpl struct {
	accountId int64
}

func (o *PlotCalculatorImpl) CalculateHit1(accountId int64, nonce int, genSig []byte, scoop, pocVersion int) *big.Int {
	plot := NewMiningPlot(uint64(accountId), uint64(nonce))
	shabal256 := shabal.NewShabal256()
	shabal256.Write(genSig)
	plot.HashScoop(shabal256, scoop)
	hash := shabal256.Sum(nil)
	return big.NewInt(0).SetBytes([]byte{hash[7], hash[6], hash[5], hash[4], hash[3], hash[2], hash[1], hash[0]})
}
func (o *PlotCalculatorImpl) CalculateHit2(genSig, scoopData []byte) *big.Int {
	shabal256 := shabal.NewShabal256()
	shabal256.Write(genSig)
	shabal256.Write(scoopData)
	hash := shabal256.Sum(nil)
	return big.NewInt(0).SetBytes([]byte{hash[7], hash[6], hash[5], hash[4], hash[3], hash[2], hash[1], hash[0]})
}
func (o *PlotCalculatorImpl) CalculateScoop(genSig []byte, height int64) int {
	shabal256 := shabal.NewShabal256()
	shabal256.Write(genSig)
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, uint64(height))
	shabal256.Write(buff)
	sum := shabal256.Sum(nil)
	hashnum := big.NewInt(0).SetBytes(sum)
	scoopnum := hashnum.Mod(hashnum, big.NewInt(int64(SCOOPS_PER_PLOT))).Int64()
	return int(scoopnum)
}
func (o *PlotCalculatorImpl) CalculateGenerationSignature(lastGenSig []byte, lastGenId int64) []byte {
	shabal256 := shabal.NewShabal256()
	shabal256.Write(lastGenSig)
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, uint64(lastGenId))
	shabal256.Write(buff)
	return shabal256.Sum(nil)
}

//baseTarget： 难度系数
func (o *PlotCalculatorImpl) CalculateDeadline(accountId int64, nonce int, genSig []byte, scoop int, baseTarget int64, pocVersion int) *big.Int {
	hit := o.CalculateHit1(accountId, nonce, genSig, scoop, pocVersion)
	return hit.Div(hit, big.NewInt(baseTarget))
}
