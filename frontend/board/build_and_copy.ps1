yarn build

# /static - for build embed
remove-item -r -force ../../static/html/board/_app
copy-item -r -force ./build/* ../../static/html/board

# /bin/static - for run, test binary
# copy-item -r -force ./build/* ../../bin/static/html/board
