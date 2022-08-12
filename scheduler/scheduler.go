package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/medialive"
	"gopkg.in/yaml.v2"
)

func main() {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	svc := medialive.New(sess)

	result, err := svc.ListChannels(nil)
	if err != nil {
		exitErrorf("Unable to list eml channels, %v", err)
	}

	fmt.Println("encoders:")

	for _, b := range result.Channels {

		fmt.Printf("** %s ID: %s **\n",
			aws.StringValue(b.Name), aws.StringValue(b.Id))
		fmt.Println()
		for _, y := range b.InputAttachments {
			fmt.Printf("- %s\n", *y.InputAttachmentName)
		}
		fmt.Println()
	}
	fmt.Println("--------------------------")
	fmt.Println()

	var SegmentationEventId, SegmentationUpidType, SegmentationDuration, SegmentationTypeId, SegmentNum, SubSegmentsExpected, SegmentsExpected, SubSegmentNum int64
	var encoderID1 string
	var scteName, time string
	var SegmentationCancelIndicator, SegmentationUpid, ReferenceActionName, FollowPoint string
	var response string

	type schedule_scte struct {
		SegmentationEventId         int64  `yaml:"segmentationEventId"`
		SegmentationCancelIndicator string `yaml:"segmentationCancelIndicator"`
		SegmentationDuration        int64  `yaml:"segmentationDuration"`
		SegmentationUpidType        int64  `yaml:"segmentationUpidType"`
		SegmentationUpid            string `yaml:"segmentationUpid"`
		SegmentationTypeId          int64  `yaml:"segmentationTypeId"`
		SegmentNum                  int64  `yaml:"segmentNum"`
		SegmentsExpected            int64  `yaml:"segmentsExpected"`
		SubSegmentNum               int64  `yaml:"subSegmentNum"`
		SubSegmentsExpected         int64  `yaml:"subSegmentsExpected"`
	}

	type schedule_input struct {
		InputAttachmentNameReference string `yaml:"inputAttachmentNameReference"`
	}

	var encoderID, actionName, inputName, startType string
	FollowPoint = "END"

	scheduleActionsImmediate := []*medialive.ScheduleAction{
		{
			ActionName: &actionName,
			ScheduleActionSettings: &medialive.ScheduleActionSettings{
				InputSwitchSettings: &medialive.InputSwitchScheduleActionSettings{
					InputAttachmentNameReference: &inputName,
				},
			},
			ScheduleActionStartSettings: &medialive.ScheduleActionStartSettings{
				ImmediateModeScheduleActionStartSettings: &medialive.ImmediateModeScheduleActionStartSettings{},
			},
		},
	}
	scheduleActionsFollow := []*medialive.ScheduleAction{
		{
			ActionName: &actionName,
			ScheduleActionSettings: &medialive.ScheduleActionSettings{
				InputSwitchSettings: &medialive.InputSwitchScheduleActionSettings{
					InputAttachmentNameReference: &inputName,
				},
			},
			ScheduleActionStartSettings: &medialive.ScheduleActionStartSettings{
				FollowModeScheduleActionStartSettings: &medialive.FollowModeScheduleActionStartSettings{
					FollowPoint:         &FollowPoint,
					ReferenceActionName: &ReferenceActionName,
				},
			},
		},
	}
	scheduleActionsFixed := []*medialive.ScheduleAction{
		{
			ActionName: &actionName,
			ScheduleActionSettings: &medialive.ScheduleActionSettings{
				InputSwitchSettings: &medialive.InputSwitchScheduleActionSettings{
					InputAttachmentNameReference: &inputName,
				},
			},
			ScheduleActionStartSettings: &medialive.ScheduleActionStartSettings{
				FixedModeScheduleActionStartSettings: &medialive.FixedModeScheduleActionStartSettings{
					Time: &time,
				},
			},
		},
	}

	params_Im := &medialive.BatchUpdateScheduleInput{
		ChannelId: &encoderID,
		Creates: &medialive.BatchScheduleActionCreateRequest{
			ScheduleActions: scheduleActionsImmediate,
		},
	}
	params_follow := &medialive.BatchUpdateScheduleInput{
		ChannelId: &encoderID,
		Creates: &medialive.BatchScheduleActionCreateRequest{
			ScheduleActions: scheduleActionsFollow,
		},
	}
	params_fixed := &medialive.BatchUpdateScheduleInput{
		ChannelId: &encoderID,
		Creates: &medialive.BatchScheduleActionCreateRequest{
			ScheduleActions: scheduleActionsFixed,
		},
	}
	SegmentationEventId = 1
	SegmentationCancelIndicator = "SEGMENTATION_EVENT_NOT_CANCELED"
	SegmentationDuration = 0
	SegmentationUpidType = 1
	SegmentNum = 0
	SegmentsExpected = 0
	SubSegmentNum = 0
	SubSegmentsExpected = 0

	scte_settings := medialive.Scte35DescriptorSettings{
		SegmentationDescriptorScte35DescriptorSettings: &medialive.Scte35SegmentationDescriptor{
			SegmentationEventId:         &SegmentationEventId,
			SegmentationCancelIndicator: &SegmentationCancelIndicator,
			SegmentationDuration:        &SegmentationDuration,
			SegmentationUpidType:        &SegmentationUpidType,
			SegmentationUpid:            &SegmentationUpid,
			SegmentationTypeId:          &SegmentationTypeId,
			SegmentNum:                  &SegmentNum,
			SegmentsExpected:            &SegmentsExpected,
			SubSegmentNum:               &SubSegmentNum,
			SubSegmentsExpected:         &SubSegmentsExpected,
		},
	}

	scheduleSCTE := []*medialive.ScheduleAction{
		{
			ActionName: &scteName,
			ScheduleActionSettings: &medialive.ScheduleActionSettings{
				Scte35TimeSignalSettings: &medialive.Scte35TimeSignalScheduleActionSettings{
					Scte35Descriptors: []*medialive.Scte35Descriptor{
						{
							Scte35DescriptorSettings: &scte_settings,
						},
					},
				},
			},
			ScheduleActionStartSettings: &medialive.ScheduleActionStartSettings{
				FollowModeScheduleActionStartSettings: &medialive.FollowModeScheduleActionStartSettings{
					ReferenceActionName: &ReferenceActionName,
					FollowPoint:         &FollowPoint,
				},
			},
		},
	}

	params1 := &medialive.BatchUpdateScheduleInput{
		ChannelId: &encoderID1,
		Creates: &medialive.BatchScheduleActionCreateRequest{
			ScheduleActions: scheduleSCTE,
		},
	}
	scheduleSCTEfixed := []*medialive.ScheduleAction{
		{
			ActionName: &scteName,
			ScheduleActionSettings: &medialive.ScheduleActionSettings{
				Scte35TimeSignalSettings: &medialive.Scte35TimeSignalScheduleActionSettings{
					Scte35Descriptors: []*medialive.Scte35Descriptor{
						{
							Scte35DescriptorSettings: &scte_settings,
						},
					},
				},
			},
			ScheduleActionStartSettings: &medialive.ScheduleActionStartSettings{
				FixedModeScheduleActionStartSettings: &medialive.FixedModeScheduleActionStartSettings{
					Time: &time,
				},
			},
		},
	}
	params2 := &medialive.BatchUpdateScheduleInput{
		ChannelId: &encoderID1,
		Creates: &medialive.BatchScheduleActionCreateRequest{
			ScheduleActions: scheduleSCTEfixed,
		},
	}

	fmt.Println("please type 'scte' for Scte35TimeSignal\n\nor type 'input' for InputSwitch: ")
	fmt.Scanln(&response)
	for {
		if response == "scte" {
			fmt.Println("Enter encoder ID: ")
			fmt.Scanln(&encoderID1)
			fmt.Println("Enter action name: ")
			fmt.Scanln(&scteName)
			fmt.Println("Enter SegmentationUpid (hexadecimal - must conatin even number of hex characters): ")
			fmt.Scanln(&SegmentationUpid)
			fmt.Println("Enter SegmentationTypeId: ")
			fmt.Scanln(&SegmentationTypeId)
			fmt.Println("Start type: 'fixed' or 'follow' ?")
			fmt.Scanln(&startType)

			for {
				if startType != "fixed" && startType != "follow" {
					fmt.Println("try again")
					fmt.Scanln(&startType)

				} else {
					break
				}
			}

			if startType == "fixed" {
				fmt.Println("Enter time for fixed start  (Format: 0000-00-00T00:00:00.000Z)\n e.g. 2022-08-21T09:30:00.000Z : ")
				fmt.Scanln(&time)
				req, resp := svc.BatchUpdateScheduleRequest(params2)
				y, err := yaml.Marshal(scheduleSCTEfixed)
				if err != nil {
					exitErrorf("error: %v", err)
				}

				err2 := req.Send()
				if err2 != nil {
					fmt.Println("\nError: please check for typos\n")
					fmt.Println(err2)
				} else {
					fmt.Println(string(y))
					fmt.Println(resp)
				}
			} else if startType == "follow" {
				fmt.Println("Enter reference action name for follow: ")
				fmt.Scanln(&ReferenceActionName)

				req, resp := svc.BatchUpdateScheduleRequest(params1)
				y, err := yaml.Marshal(scheduleSCTE)
				if err != nil {
					exitErrorf("error: %v", err)
				}

				err2 := req.Send()
				if err2 != nil {
					fmt.Println("\nError: please check for typos\n")
					fmt.Println(err2)
				} else {
					fmt.Println(string(y))
					fmt.Println(resp)
				}
			}

			fmt.Println("schedule again? (y/n): ")
			fmt.Scanln(&response)
			if response == "y" {
				fmt.Println("please type 'scte' for Scte35TimeSignal\n\nor type 'input' for InputSwitch: ")
				fmt.Scanln(&response)
				continue
			} else {
				break
			}

		} else if response == "input" {
			fmt.Println("Enter encoder ID: ")
			fmt.Scanln(&encoderID)
			fmt.Println("Enter action name: ")
			fmt.Scanln(&actionName)
			fmt.Println("Enter input name: ")
			fmt.Scanln(&inputName)
			fmt.Println("Start type: 'fixed', 'follow', or 'immediate' ?")
			fmt.Scanln(&startType)
			for {
				if startType != "fixed" && startType != "follow" && startType != "immediate" {
					fmt.Println("try again")
					fmt.Scanln(&startType)

				} else {
					break
				}
			}

			if startType == "immediate" {
				req, resp := svc.BatchUpdateScheduleRequest(params_Im)
				y, err := yaml.Marshal(scheduleActionsImmediate)
				if err != nil {
					exitErrorf("error: %v", err)
				}
				err2 := req.Send()
				if err2 != nil {
					fmt.Println("\nError: please check for typos\n")
					fmt.Println(err2)
				} else {
					fmt.Println(string(y))
					fmt.Println(resp)
				}
			} else if startType == "follow" {
				fmt.Println("Enter reference action name for follow: ")
				fmt.Scanln(&ReferenceActionName)
				req, resp := svc.BatchUpdateScheduleRequest(params_follow)
				y, err := yaml.Marshal(scheduleActionsFollow)
				if err != nil {
					exitErrorf("error: %v", err)
				}

				err2 := req.Send()
				if err2 != nil {
					fmt.Println("\nError: please check for typos\n")
					fmt.Println(err2)
				} else {
					fmt.Println(string(y))
					fmt.Println(resp)
				}
			} else if startType == "fixed" {
				fmt.Println("Enter time for fixed start (Format: 0000-00-00T00:00:00.000Z)\n e.g. 2022-08-21T09:30:00.000Z")
				fmt.Scanln(&time)
				req, resp := svc.BatchUpdateScheduleRequest(params_fixed)
				y, err := yaml.Marshal(scheduleActionsFixed)
				if err != nil {
					exitErrorf("error: %v", err)
				}
				err2 := req.Send()
				if err2 != nil {
					fmt.Println("\nError: please check for typos\n")
					fmt.Println(err2)
				} else {
					fmt.Println(string(y))
					fmt.Println(resp)
				}
			}

			fmt.Println("schedule again? (y/n): ")
			fmt.Scanln(&response)
			if response == "y" {
				fmt.Println("please type 'scte' for Scte35TimeSignal\n\nor type 'input' for InputSwitch: ")
				fmt.Scanln(&response)
				continue
			} else {
				break
			}
		} else if response != "scte" && response != "input" {
			fmt.Println("try again, type 'scte' for Scte35TimeSignal\n\nor 'input' for InputSwitch")
			fmt.Scanln(&response)
			continue
		}
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
