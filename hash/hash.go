package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
)

func SHAtoBigNum(sha string) (*big.Int, error) {
	if len(sha) != 64 {
		return nil, errors.New("sha must be 64 length")
	}
	big_sum := big.NewInt(0)
	for i := 0; i < 64; i += 8 {
		i64, err := strconv.ParseInt(sha[i:i+8], 16, 64)
		if err != nil {
			log.Print(err)
		}

		big_i64 := big.NewInt(i64)
		// 10
		jousyuu := big.NewInt(int64(16))
		// i big int
		jousuu_i := big.NewInt(int64(56 - i))
		// 1000000を追加するためのもの
		jousuu_konnkai := jousyuu.Exp(jousyuu, jousuu_i, nil)

		add_big_i64_num := big_i64.Mul(
			big_i64,
			jousuu_konnkai,
		)
		big_sum.Add(big_sum, add_big_i64_num)
	}
	return big_sum, nil
}

func SHA256Threshold(hp int64, diff int64) *big.Int {
	// ( 2 ** 256) * hp / diff
	big2 := big.NewInt(2)
	big256 := big.NewInt(256)
	bighp := big.NewInt(hp)
	bigdiff := big.NewInt(diff)

	threshold := big256.Exp(big2, big256, nil).Mul(big256, bighp).Div(big256, bigdiff)
	return threshold
}

func SHA256(blockword string) []byte {
	// SHA256(blockword)
	hasher := sha256.New()
	hasher.Write([]byte(blockword))
	return hasher.Sum(nil)
}

func IsOKHash(hp int64, diff int64, blockword string) bool {
	// SHA256(blockword) <= 2 ** 256 * hp / diff

	threshold := SHA256Threshold(hp, diff)
	exsha := SHA256(blockword)
	// fmt.Println(exsha)
	result_exsha, err := SHAtoBigNum(hex.EncodeToString(exsha))
	if err != nil {
		log.Print(err)
	}

	log.Print("isOKHash", fmt.Sprint(threshold))
	log.Print("isOKHash: ", result_exsha)
	log.Print("isOKHash: ", threshold.Cmp(result_exsha))
	log.Print("isOKHash: ", 0 <= threshold.Cmp(result_exsha))

	return 0 <= threshold.Cmp(result_exsha)
}
