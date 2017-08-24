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
	int=0
	while [ $int -lt 100 ]; do
		int=$(($int + 2))
		echo "a" | dialog --title "D-voting" --gauge "$title" 0 0 $int
		sleep 0.02
	done
	echo "a" | dialog --title "D-voting" --msgbox "$title" 0 0 
	sleep 0.5
}

main()
{
	#main menu
	dialog --no-cancel --title "Menu" --menu "" 200 100 6\
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
	sleep 0.5
	main
}

down()
{
	if [ "$(docker ps | grep chaincode)" == "" ]; then
		dialog --title "Network is down already" --infobox "Network is down already: aborting" 200 100
		sleep 2
		main
	fi
	dialog --title "Process: bring down the network" --infobox "bringing down the network..." 200 100
	service docker restart
	main
}

chaincode()
{
	chaincode_build()
	{
		codename="voting"
		docker exec -d chaincode bash -c "cd $codename; go build 2> /dev/null;"
		dialog --title "Chaincode building" --infobox "Chaincode is building" 200 100
		sleep 2
	}
	chaincode_run()
	{
		codename="voting"
		ping=$(docker exec chaincode ps aux | grep $codename)
		if [ "$ping" != ""  ]; then
			dialog --title "Chaincode is running" --infobox "Chaincode is already running, aborting..." 200 100
			sleep 2
		fi
		docker exec -d chaincode bash -c "cd $codename;CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc:0 ./$codename 2> /dev/null"
	}
	chaincode_kill()
	{
		codename="voting"
		ping=$(docker exec chaincode ps aux | grep $codename)
		if [ "$ping" == ""  ]; then
			dialog --title "Chaincode is not running" --infobox "Chaincode is not running, aborting..." 200 100
			sleep 2
		fi
		docker exec -d chaincode bash -c "killall $codename"
		dialog --title "process killed" --infobox "Chaincode killed" 200 100
	}
	chaincode_install()
	{
		codename="voting"
		docker exec -d cli bash -c "peer chaincode install -p chaincodedev/chaincode/$codename -n mycc -v 0"
	}
	chaincode_instantiate()
	{
		docker exec -d cli bash -c "peer chaincode instantiate -n mycc -v 0 -c '{\"Args\":[\"a\",\"10\"]}' -C myc"
	}
	chaincode_invoke()
	{
		dialog --title "" --infobox "" 200 100
		dialog --title "Not yet" --infobox "Not yet" 200 100
		sleep 2
	}
	chaincode_exit()
	{
		main
	}
	dialog --title "Chaincode" --menu "" 200 100 7\
		build "Build chaincode"\
		run "run chaincode"\
		kill "kill running chaincode"\
		install "install chaincode"\
		instantiate "instantate chaincode"\
		invoke "invoke chaincode"\
		exit "exit the system" 2> /tmp/netup.chaincode
	result=$?
	if [ $result -eq 1 ]; then
		main
	elif [ $result -eq 255 ]; then
		exit 255
	fi
	ans=$(cat /tmp/netup.chaincode)
	rm /tmp/netup.chaincode
	chaincode_$ans
	chaincode
}

login()
{
	login_username()
	{
		username=""
		while [ "$username" == "" ]; do
			dialog --title "Username" --inputbox "Please input username" 200 100 2> /tmp/netup.login
			result=$?
			if [ $result -eq 1 ]; then
				main
			elif [ $result -eq 255 ]; then
				exit 255
			fi
			username=$(cat /tmp/netup.login)
			rm /tmp/netup.login
		done
	}
	login_password()
	{
		password=""
		while [ "$password" == "" ]; do
			dialog --insecure --title "Password" --passwordbox "Please input password" 200 100 2> /tmp/netup.login
			result=$?
			if [ $result -eq 1 ]; then
				main
			elif [ $result -eq 255 ]; then
				exit 255
			fi
			password=$(cat /tmp/netup.login)
			rm /tmp/netup.login
		done
	}
	login_username
	login_password
	login_adduser()
	{
		login_adduser_username()
		{
			newusername=""
			while [ "$newusername" == "" ]; do
				dialog --title "Adduser" --inputbox "Please input new username" 200 100 2> /tmp/netup.login.adduser
				result=$?
				if [ $result -eq 1 ]; then
					main
				elif [ $result -eq 255 ]; then
					exit 255
				fi
				newusername=$(cat /tmp/netup.login.adduser)
				rm /tmp/netup.login.adduser
			done
		}
		login_adduser_password()
		{
			newpassword=""
			while [ "$newpassword" == "" ]; do
				dialog --insecure --title "Adduser" --passwordbox "Please input the password" 200 100 2> /tmp/netup.login.adduser
				result=$?
				if [ $result -eq 1 ]; then
					main
				elif [ $result -eq 255 ]; then
					exit 255
				fi
				newpassword=$(cat /tmp/netup.login.adduser)
				rm /tmp/netup.login.adduser
			done
		}
		login_adduser_password_2()
		{
			newpassword2=""
			while [ "$newpassword2" == "" ]; do
				dialog --insecure --title "Adduser" --passwordbox "Please retype the password" 200 100 2> /tmp/netup.login.adduser
				result=$?
				if [ $result -eq 1 ]; then
					main
				elif [ $result -eq 255 ]; then
					exit 255
				fi
				newpassword2=$(cat /tmp/netup.login.adduser)
				rm /tmp/netup.login.adduser
			done
		}
		login_adduser_username
		login_adduser_password
		login_adduser_password_2
		result=$(docker exec cli bash -c "peer chaincode query -n mycc -c '{\"Args\":[\"verify\", \"$newusername\", \"$newpassword\"]}' -C myc --logging-level=error")
		if [ "$result" != "Query Result: NO SUCH USER" ]; then
			dialog --title "User exists" --msgbox "User exist: use another name" 200 100
			sleep 2
			main
		elif [ "$newpassword" != "$newpassword2" ]; then
			dialog --title "Password mismatch" --msgbox "Password mismatch, please retry" 200 100
			sleep 2
			main
		else 
			docker exec cli bash -c "peer chaincode invoke -n mycc -c '{\"Args\":[\"adduser\", \"$newusername\", \"$newpassword\"]}' -C myc --logging-level=error"
		fi
		main
	}
	if [ "$username" == "adduser" ]; then
		if [ "$password" == "iamroot" ]; then
			login_adduser
		else 
			dialog --title "Wrong password" --msgbox "Wrong password..." 200 100
			login
		fi
	fi
	result=$(docker exec cli bash -c "peer chaincode query -n mycc -c '{\"Args\":[\"verify\", \"$username\", \"$password\"]}' -C myc --logging-level=error")
	if [ "$result" == "Query Result: NO SUCH USER" ]; then
		dialog --title "No such user" --msgbox "User not found..." 200 100
		login
	fi
	if [ "$result" == "Query Result: PASSWORD WRONG" ]; then
		dialog --title "Wrong password" --msgbox "Wrong password..." 200 100
		login
	fi
	if [ "$result" == "Query Result: VOTE USED" ]; then
		dialog --title "Account status" --msgbox "Vote used..." 200 100
		login
	fi
	if [ "$result" == "Query Result: VOTE NOT USED" ]; then
		dialog --title "Account status" --msgbox "Vote not used..." 200 100
		dialog --title "Vote for" --menu "" 200 100 4\
			1 "CANDIDATE_1"\
			2 "CANDIDATE_2"\
			3 "CANDIDATE_3"\
			exit "go to the last page" 2> /tmp/netup.chaincode
		result=$?
		if [ $result -eq 1 ]; then
			main
		elif [ $result -eq 255 ]; then
			exit 255
		fi
		ans=$(cat /tmp/netup.chaincode)
		rm /tmp/netup.chaincode
		if [ $ans == "exit" ]; then
			main
		fi
		result=$(docker exec cli bash -c "peer chaincode invoke -n mycc -c '{\"Args\":[\"votefor\", \"$username\", \"$password\", \"CANDIDATE_$ans\"]}' -C myc --logging-level=error")
		sleep 2
		ans1=$(docker exec cli bash -c "peer chaincode query -n mycc -c '{\"Args\":[\"seepoll\", \"CANDIDATE_1\"]}' -C myc --logging-level=error")
		ans2=$(docker exec cli bash -c "peer chaincode query -n mycc -c '{\"Args\":[\"seepoll\", \"CANDIDATE_2\"]}' -C myc --logging-level=error")
		ans3=$(docker exec cli bash -c "peer chaincode query -n mycc -c '{\"Args\":[\"seepoll\", \"CANDIDATE_3\"]}' -C myc --logging-level=error")
		dialog --title "Vote succeeded" --msgbox "Now poll:\n Candidate 1: $ans1\n Candidate 2: $ans2\n Candidate 3: $ans3\n" 200 100
		main
	fi
}

begin
main
