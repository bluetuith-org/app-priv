package tlog

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/spf13/cast"
)

const timeFmt = "2006-01-02 15:04:05"

const (
	segSep    = ", "
	kvSep     = "="
	headerSep = ": "
)

func getAttrBuf(r slog.Record) *Buffer {
	by := NewBuffer()
	by.Reset()

	*by = append(*by, r.Message...)
	*by = append(*by, headerSep...)

	appendSep := false
	r.Attrs(func(a slog.Attr) bool {
		a.Value = a.Value.Resolve()

		if appendSep {
			*by = append(*by, segSep...)
		}

		*by = append(*by, a.Key...)
		*by = append(*by, kvSep...)

		switch a.Value.Kind() {
		case slog.KindBool:
			*by = strconv.AppendBool(*by, a.Value.Bool())

		case slog.KindInt64:
			*by = strconv.AppendInt(*by, a.Value.Int64(), 10)

		case slog.KindUint64:
			*by = strconv.AppendUint(*by, a.Value.Uint64(), 10)

		case slog.KindDuration:
			var arr [32]byte
			n := formatDuration(a.Value.Duration(), &arr)
			*by = append(*by, arr[n:]...)

		case slog.KindTime:
			*by = a.Value.Time().AppendFormat(*by, timeFmt)

		case slog.KindFloat64:
			*by = strconv.AppendFloat(*by, a.Value.Float64(), 'f', -1, 64)

		case slog.KindString:
			*by = append(*by, a.Value.String()...)

		case slog.KindAny:
			*by = append(*by, cast.ToString(a.Value.Any())...)
		}

		appendSep = true

		return true
	})

	return by
}

// formatDuration formats the representation of d into the end of buf and
// returns the offset of the first character.
func formatDuration(d time.Duration, buf *[32]byte) int {
	// Largest time is 2540400h10m10.000000000s
	w := len(buf)

	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	if u < uint64(time.Second) {
		// Special case: if duration is smaller than a second,
		// use smaller units, like 1.2ms
		var prec int
		w--
		buf[w] = 's'
		w--
		switch {
		case u == 0:
			buf[w] = '0'
			return w
		case u < uint64(time.Microsecond):
			// print nanoseconds
			prec = 0
			buf[w] = 'n'
		case u < uint64(time.Millisecond):
			// print microseconds
			prec = 3
			// U+00B5 'µ' micro sign == 0xC2 0xB5
			w-- // Need room for two bytes.
			copy(buf[w:], "µ")
		default:
			// print milliseconds
			prec = 6
			buf[w] = 'm'
		}
		w, u = fmtFrac(buf[:w], u, prec)
		w = fmtInt(buf[:w], u)
	} else {
		w--
		buf[w] = 's'

		w, u = fmtFrac(buf[:w], u, 9)

		// u is now integer seconds
		w = fmtInt(buf[:w], u%60)
		u /= 60

		// u is now integer minutes
		if u > 0 {
			w--
			buf[w] = 'm'
			w = fmtInt(buf[:w], u%60)
			u /= 60

			// u is now integer hours
			// Stop at hours because days can be different lengths.
			if u > 0 {
				w--
				buf[w] = 'h'
				w = fmtInt(buf[:w], u)
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}

	return w
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros. It omits the decimal
// point too when the fraction is 0. It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	pr := false
	for range prec {
		digit := v % 10
		pr = pr || digit != 0
		if pr {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if pr {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}
