#!/bin/sh
alias dialog="dialog --ascii-lines"


title='\n
    ____                 __  _            \n
   / __ \   _   ______  / /_(_)___  ____ _\n
  / / / /  | | / / __ \/ __/ / __ \/ __ `/\n
 / /_/ /   | |/ / /_/ / /_/ / / / / /_/ / \n
/_____/    |___/\____/\__/_/_/ /_/\__, /  \n
                                 /____/   \n'

#dialog --title "Logging mode"
#docker-compose -f docker-compose-simple.yaml up

begin()
{
	#check root
	if [ "$(whoami)" != "root" ]; then
		dialog --title "Permission denied" --infobox "Please execute as root" 200 100
		exit
	fi
	
	#welcoming banner
	cat /dev/null > /tmp/netup.main
	int=0
	while [ $int -lt 100 ]; do
		int=$(($int + 2))
		echo "a" | dialog --title "D-voting" --gauge "$title" 0 0 $int
		sleep 0.02
	done
	echo "a" | dialog --title "D-voting" --gauge "$title" 0 0 $int
	echo "a" | dialog --title "D-voting" --msgbox "$title" 0 0 
	sleep 0.5
}

main()
{
	#main menu
	dialog --title "Menu" --menu "" 200 100 6\
		up "bring the network up"\
		down "bring the network down"\
		prune "kill all unusing container"\
		chaincode "perform operation on chaincode"\
		login "login to the d-voting system"\
		exit "exit the system" 2> /tmp/netup.main
	result=$?
	if [ $result -eq 1 ]; then
		exit 1
	elif [ $result -eq 255 ]; then
		exit 255
	fi
	ans=$(cat /tmp/netup.main)
	rm /tmp/netup.main
	$ans
}

up()
{
	#check if machine is up
	if [ "$(docker ps | grep cli)" != "" ]; then
		dialog --title "Container running" --yesno "Container running, kill all and restart? " 200 100
		result=$?
		if [ $result -eq 1 ]; then
			main
		elif [ $result -eq 255 ]; then
			exit 255
		fi
		dialog --title "Process: restart network" --infobox "restarting network..." 200 100
		service docker restart
		yes | docker container prune
	fi

	if [ "$(docker container ls | grep cli)" != "" ]; then
		dialog --title "Container corpse discovered" --yesno "Container corpse discovered, cleanup? " 200 100
		result=$?
		if [ $result -eq 1 ]; then
			main
		elif [ $result -eq 255 ]; then
			exit 255
		fi
		dialog --title "Process: restart network" --infobox "restarting network..." 200 100
		yes | docker container prune
	fi
	#bring up the network
	dialog --title "Process: bring up the network" --infobox "bringing up the network..." 200 100
	docker-compose -f docker-compose-simple.yaml up -d 2>&1 > /dev/null
	echo "a" | dialog --title "D-voting" --gauge "Checking chaincode" 200 100 25
	sleep 1
	echo "a" | dialog --title "D-voting" --gauge "Checking cli" 200 100 50
	sleep 1.5
	echo "a" | dialog --title "D-voting" --gauge "Checking peer" 200 100 75
	sleep 0.5
	echo "a" | dialog --title "D-voting" --gauge "Checking orderer" 200 100 100
	sleep 0.5
	main
}

prune()
{
	dialog --title "Process: pruning" --infobox "pruning..." 200 100
	yes | docker container prune 2>&1 > /dev/null
	main
}

down()
{
	if [ "$(docker ps | grep chaincode)" == "" ]; then
		dialog --title "Network is down already" --infobox "Network is down already: aborting" 200 100
		sleep 1
		main
	fi
	dialog --title "Process: bring down the network" --infobox "bringing down the network..." 200 100
	service docker restart
	main
}

chaincode()
{
	echo chaincode
}

login()
{
	echo login
}

begin
main
