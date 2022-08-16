#!/usr/bin/sh

dir=dir
file=template.monit
link=linka
mkdir -p ${dir}

if [ -L "${dir}" ]; then
    echo "${dir} is link"
elif [ -d "${dir}" ]; then 
    echo "${dir} is a directory"
else 
    echo "${dir} is a regular file"
fi

ln -sf ${file} ${dir} ## error
ln -sf ${dir} ${link}
if [ -L "${link}" ]; then
    echo "${link} is link"
elif [ -d "${link}" ]; then 
    echo "${link} is a directory"
else 
    echo "${link} is a regular file"
fi

if [ ! -L "./${file}" ]; then
    echo "${file} is not link"
fi

if [ ! -L "./${link}" ]; then
    echo "${link} is not link"
fi