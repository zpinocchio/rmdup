# rmdup
this script is change from seqkit rmdup -by_seq, same algorithm but input R1 R2,output R1 R2<br>
for more information see [shenwei356/seqkit](https://github.com/zpinocchio/rmdup/upload/master)<br>

# Usage


# Core algorithm
```go
subjectR1 = xxhash.Sum64(recordR1.Seq.Seq)
subjectR2 = xxhash.Sum64(recordR2.Seq.Seq)
_, ok := counterR1[subjectR1]
_,ok2 := counterR2[subjectR2]
if ok || ok2{ // duplicated
  counterR1[subjectR1]++
	counterR2[subjectR2]++
	removed++
  } else { // new one
	recordR1.FormatToWriter(outfhR1, LineWidth)
	recordR2.FormatToWriter(outfhR2, LineWidth)
	counterR1[subjectR1]++
	counterR2[subjectR2]++
	}
```
