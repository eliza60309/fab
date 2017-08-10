#!/bin/sh

error_message()
{
	echo "#######################################"
	echo "error occured while performing netstart"
	echo "#######################################"
}

if [ $? -eq 1 ] ;
then
       error_message	
fi

