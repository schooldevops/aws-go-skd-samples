package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	PROFILE = "sdk-user"
	REGION  = "ap-northeast-2"
)

var (
	SecretName  = "myproject/devops/configs/v1"
	Description = "This secret is my Project Secret"
	Secrets     = map[string]string{
		"myDBName":     "dbname",
		"myDBPassword": "dbpassword",
		"dbPort":       "dbPorts",
	}
	UpdateSecrets = map[string]string{
		"myDBName":     "modDbname",
		"myDBPassword": "modDbpassword",
		"dbPort":       "modDbPorts",
	}
	PutSecrets = map[string]string{
		"myDBName":     "putDbname",
		"myDBPassword": "putDbpassword",
		"dbPort":       "putDbPorts",
	}
)

// getSecretManager()
// secret manager 에 접근하기 위한 세션을 생성한다.
func getSecretManager() *secretsmanager.SecretsManager {
	// Secret Manager 에 접근할 세션을 생성한다.
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Profile: "sdk-user",
			Config: aws.Config{
				Region: aws.String("ap-northeast-2"),
			},
		},
	)

	if err != nil {
		log.Fatalf("Cannot create session %v", err)
	}

	smgr := secretsmanager.New(sess)

	return smgr
}

// Secret 을 신규로 생성한다.
func createSecrets(smgr *secretsmanager.SecretsManager, inputValue *secretsmanager.CreateSecretInput) *secretsmanager.CreateSecretOutput {
	result, err := smgr.CreateSecret(inputValue)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodePreconditionNotMetException:
				fmt.Println(secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

// Secret 정보를 조회한다.
func getSecretInfos(smgr *secretsmanager.SecretsManager, secretId string) *secretsmanager.DescribeSecretOutput {
	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretId),
	}

	result, err := smgr.DescribeSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

// 리소스 정책을 조회한다.
func getResourcePolicy(smgr *secretsmanager.SecretsManager, secretId string) *secretsmanager.GetResourcePolicyOutput {
	input := &secretsmanager.GetResourcePolicyInput{
		SecretId: aws.String(secretId),
	}

	result, err := smgr.GetResourcePolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

// 시크릿 키에 대한 값을 조회한다.
func getSecretValue(smgr *secretsmanager.SecretsManager, secretId string) *secretsmanager.GetSecretValueOutput {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}

	result, err := smgr.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result

}

// 시크릿 정보를 업데이트한다.
// 기본 시크릿 정보를 변경하고,
func updateSecret(smgr *secretsmanager.SecretsManager, input *secretsmanager.UpdateSecretInput) *secretsmanager.UpdateSecretOutput {

	result, err := smgr.UpdateSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodePreconditionNotMetException:
				fmt.Println(secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

// 시크릿 정보를 업데이트한다.
func putSecretValue(smgr *secretsmanager.SecretsManager, input *secretsmanager.PutSecretValueInput) *secretsmanager.PutSecretValueOutput {

	result, err := smgr.PutSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

func deleteSecret(smgr *secretsmanager.SecretsManager, secretId string, optional ...int64) *secretsmanager.DeleteSecretOutput {
	recoveryDay := int64(7)
	if len(optional) > 0 {
		recoveryDay = optional[0]
	}
	input := &secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(recoveryDay),
		SecretId:             aws.String(secretId),
	}

	result, err := smgr.DeleteSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

func main() {

	smgr := getSecretManager()

	secretString, _ := json.Marshal(Secrets)

	input := &secretsmanager.CreateSecretInput{
		Description:  aws.String(Description),
		Name:         aws.String(SecretName),
		SecretString: aws.String(string(secretString)),
	}

	result := createSecrets(smgr, input)
	fmt.Println("----- Create Secret Success: ----- ")
	fmt.Println(result)

	secretInfo := getSecretInfos(smgr, SecretName)
	fmt.Println("----- Get Secret Infos: ----- ")
	fmt.Println(secretInfo)

	secretPolicy := getResourcePolicy(smgr, SecretName)
	fmt.Println("----- Get Resource Policy: ----- ")
	fmt.Println(secretPolicy)

	myDBPassword := getSecretValue(smgr, SecretName)
	fmt.Println("----- Get My DB Password: ----- ")
	fmt.Println(myDBPassword)

	modSecretString, _ := json.Marshal(UpdateSecrets)
	updateValues := &secretsmanager.UpdateSecretInput{
		Description:  aws.String("Modified Descriptions."),
		SecretId:     aws.String(SecretName),
		SecretString: aws.String(string(modSecretString)),
	}

	updateResult := updateSecret(smgr, updateValues)
	fmt.Println("----- Update Secrets: ----- ")
	fmt.Println(updateResult)

	putSecretString, _ := json.Marshal(PutSecrets)
	putValues := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(SecretName),
		SecretString: aws.String(string(putSecretString)),
	}

	putResult := putSecretValue(smgr, putValues)
	fmt.Println("----- Put Secrets: ----- ")
	fmt.Println(putResult)

	deleteResult := deleteSecret(smgr, SecretName, 30)
	fmt.Println("----- Delete Secrets: ----- ")
	fmt.Println(deleteResult)

}
