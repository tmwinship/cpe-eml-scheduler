# cpe-eml-scheduler

`scheduler.go` generates a list of all MediaLive channels contained in `us-west-2` along with the corresponding `ID` and the active inputs attached to each channel

You are given the option to schedule a `Scte-35 Time Signal` by typing `scte`

or schedule an `Input Switch` by typing `input`


You will be prompted to enter the encoder `ID` associated with the channel you wish to schedule for. Use the generated list to find the correct ID

you will then be prompted to create an `action name` for the schedule action you are currently creating

### For `scte` option: ###

enter the `Segmentation UPID`: this must be a string in hexadecimal form. make sure to have an even number of hex characters, or else an error will occur

enter the `Segmentation Type ID` refer to this document for different type ids: https://wagtail-prod-storage.s3.amazonaws.com/documents/SCTE_35_2022.pdf pages 66-68

You will then be given 2 start type options: Fixed or Follow

#### a scte35 can only follow input switches, they cannot follow other scte msgs or else error will occur ####

`fixed` time format: 0000-00-00T00:00:00.000Z 

e.g. 2022-08-21T09:30:00.000Z

for `follow` start type: enter the `action name` of the schedule action you would like to follow

### for `input` option: ###

You will be prompted to enter an input name. See generated list to find all input names attached to each channel

You will then be given 3 start type options: Fixed, Follow, or Immediate

`fixed` time format: 0000-00-00T00:00:00.000Z 

e.g. 2022-08-21T09:30:00.000Z

for `follow` start type: enter the `action name` of the schedule action you would like to follow

pick `immediate` to schedule the input switch immediately
