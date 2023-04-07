for top in pkg internal/pkg
do
    for d in $(find $top -type d)
    do
        if [ ! -f $d/doc.go ]; then
            if ls $d/*.go > /dev/null 2>&1; then
                echo $d/doc.go
                echo "package $(basename $d) // import \"github.com/changaolee/skeleton/$d\"" > $d/doc.go
            fi
        fi
    done
done
