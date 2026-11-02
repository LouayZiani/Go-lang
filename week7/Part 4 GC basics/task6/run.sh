set -e
go build -o /tmp/gcdemo ./task6/gcdemo
for g in 50 100 200 400; do
  GOGC=$g /tmp/gcdemo
done
GOGC=off /tmp/gcdemo
GOGC=off GOMEMLIMIT=100MiB /tmp/gcdemo
echo
echo "---- gctrace (first 5 lines + last one) ----"
GODEBUG=gctrace=1 /tmp/gcdemo 2>&1 | grep '^gc ' | sed -n '1,5p;$p'
