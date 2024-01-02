package random

import (
	"math"
	"math/rand"
	"time"
	"unicode/utf8"
)

type Characters = string

var (
	randomizer = New() // default randomizer
)

const (
	CharsChineseNumeralsUpper Characters = "零壹貳參肆伍陸柒捌玖拾"
	CharsChineseNumeralsLower Characters = "〇一二三四五六七八九十"
	CharsChineseNumerals                 = CharsChineseNumeralsLower + CharsChineseNumeralsUpper
	CharsGreekLetters         Characters = "αβγδεζηθικλμνξοπρστυφχψω"
	CharsRomanLettersLower    Characters = "abcdefghijklmnopqrstuvwxyz"
	CharsRomanLettersUpper    Characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsRomanNumerals        Characters = "ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫ"
	CharsRomanAll                        = CharsRomanLetters + CharsRomanNumerals
	CharsArabicNumerals       Characters = "0123456789"
	Symbols                   Characters = "`" + `~!@#$%^&*()-_+={}[]|\;:"<>,./?`
	CharsRomanLetters                    = CharsRomanLettersLower + CharsRomanLettersUpper
	CharsDefault                         = CharsArabicNumerals + CharsRomanLetters
	CharsHexDigits                       = CharsArabicNumerals + "abcdef"
	CharsSafeURL                         = CharsArabicNumerals + CharsRomanLettersLower + CharsRomanLettersUpper + "-._~"
	CharsAll                             = CharsChineseNumerals + CharsGreekLetters + CharsRomanAll + CharsDefault
	CharsAll2                            = CharsChineseNumerals + CharsGreekLetters + CharsRomanAll + CharsDefault + Symbols

	Float32ofMaxInt32 = float32(math.MaxInt32)
	Float64ofMaxInt64 = float64(math.MaxInt64)
)

type Config struct {
	Chars Characters
	Seed  int64
}

type Random struct {
	core     *rand.Rand
	seed     int64
	chars    []rune
	charsLen int64
}

func (r *Random) Chars() Characters {

	return string(append(make([]rune, 0, r.charsLen), r.chars...))
}

func (r *Random) RuneAtChars(n int) rune {
	return r.chars[n]
}

func (r *Random) SetChars(chars Characters) {
	r.chars = []rune(chars)
	r.charsLen = int64(len(r.chars))
}

func (r *Random) Seed(seed int64) {
	r.core = rand.New(rand.NewSource(seed))
}

func (r *Random) WithChars(chars Characters) *Random {
	return &Random{
		core:     r.core,
		seed:     r.seed,
		chars:    []rune(chars),
		charsLen: int64(utf8.RuneCountInString(chars)),
	}
}

func (r *Random) WithSeed(seed int64) *Random {
	return &Random{
		core:     rand.New(rand.NewSource(seed)),
		seed:     seed,
		chars:    r.chars,
		charsLen: r.charsLen,
	}
}

func (r *Random) Bytes(size int) []byte {
	buf := make([]byte, size)
	r.core.Read(buf)

	return buf
}

func (r *Random) Rune() rune {
	return r.chars[r.core.Int63()%r.charsLen]
}

func (r *Random) Runes(size int) []rune {
	buf := make([]rune, size)

	for i := 0; i < size; i++ {
		buf[i] = r.chars[r.core.Int63()%r.charsLen]
	}

	return buf
}

func (r *Random) RandString(size int) string {
	return string(r.Runes(size))
}

func (r *Random) RandInt(min, max int) int {
	return min + (r.core.Int() % (max - min))
}

func (r *Random) Int() int {
	return r.core.Int()
}

func (r *Random) Int31() int32 {
	return r.core.Int31()
}

func (r *Random) Int63() int64 {
	return r.core.Int63()
}

func (r *Random) Uint32() uint32 {
	return r.core.Uint32()
}

func (r *Random) Uint64() uint64 {
	return r.core.Uint64()
}

func (r *Random) Float32() float32 {
	return r.core.Float32()
}

func (r *Random) FastFloat32() float32 {
	return float32(r.core.Int63()%math.MaxInt32) / Float32ofMaxInt32
}

func (r *Random) Float64() float64 {
	return r.core.Float64()
}

func (r *Random) FastFloat64() float64 {
	return float64(r.core.Int63()) / Float64ofMaxInt64
}

func New(conf ...Config) *Random {
	var (
		chars = CharsDefault
		seed  = time.Now().UnixNano()
	)

	if len(conf) > 0 {
		config := conf[0]

		if config.Chars != "" {
			chars = config.Chars
		}

		if config.Seed > 0 {
			seed = config.Seed
		}
	}

	return &Random{
		core:     rand.New(rand.NewSource(seed)),
		seed:     seed,
		chars:    []rune(chars),
		charsLen: int64(len(chars)),
	}
}

// WithChars 设置字符串
func WithChars(chars Characters) *Random {
	return randomizer.WithChars(chars)
}

// WithSeed 设置随机种子
func WithSeed(seed int64) *Random {
	return randomizer.WithSeed(seed)
}

func Bytes(size int) []byte {
	return randomizer.Bytes(size)
}

func BytesF(size int) func() []byte {
	return func() []byte {
		return randomizer.Bytes(size)
	}
}

func Rune() rune {
	return randomizer.Rune()
}

func RuneF() func() rune {
	return func() rune {
		return randomizer.Rune()
	}
}

func Runes(size int) []rune {
	return randomizer.Runes(size)
}

func RunesF(size int) func() []rune {
	return func() []rune {
		return randomizer.Runes(size)
	}
}

func RandInt(min, max int) int {
	return randomizer.RandInt(min, max)
}

func RandIntF(min, max int) func() int {
	return func() int {
		return randomizer.RandInt(min, max)
	}
}

func RandString(size int) string {
	return randomizer.RandString(size)
}

func RandStringF(size int) func() string {
	return func() string {
		return randomizer.RandString(size)
	}
}

func Int() int {
	return randomizer.Int()
}

func Int31() int32 {
	return randomizer.Int31()
}

func Int63() int64 {
	return randomizer.Int63()
}

func Uint32() uint32 {
	return randomizer.Uint32()
}

func Uint64() uint64 {
	return randomizer.Uint64()
}

func Float32() float32 {
	return randomizer.Float32()
}

func FastFloat32() float32 {
	return randomizer.FastFloat32()
}

func Float64() float64 {
	return randomizer.Float64()
}

func Choice[T any](items []T) T {
	return items[Int()%len(items)]
}
