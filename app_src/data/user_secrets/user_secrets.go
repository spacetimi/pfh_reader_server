package user_secrets

import (
    "context"
    "errors"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service/storage_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/encryption_utils"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "strconv"
)

const kBlobName = "user_secrets"

// Implements IBlob
type UserSecretsBlob struct {
    UserId int64
    UserSecrets []*UserSecret

    storage_typedefs.BlobDescriptor `bson:"ignore"`
}

type UserSecret struct {
    SecretName string
    SecretDataEncrypted string
}

func LoadByUserId(userId int64, ctx context.Context, create bool) (*UserSecretsBlob, error) {
    userSecrets := newUserSecretsBlob(userId)

    err := storage_service.GetBlobByPrimaryKeys(userSecrets, ctx)
    if err != nil {
        if !create {
            return nil, errors.New("error getting " + kBlobName + " blob: " + err.Error())
        }

        userSecrets, err = Create(userId, ctx)
        if err != nil {
            return nil, errors.New("error creating " + kBlobName + " blob: " + err.Error())
        }
    }

    return userSecrets, nil

}

func Create(userId int64, ctx context.Context) (*UserSecretsBlob, error) {
    userSecrets := newUserSecretsBlob(userId)
    err := storage_service.SetBlob(userSecrets, ctx)
    if err != nil {
        return nil, errors.New("error saving " + kBlobName + " blob: " + err.Error())
    }

    return userSecrets, nil
}

func (blob *UserSecretsBlob) GetSecret(secretName string, masterPassword string) (string, error) {
    userSecret, err := blob.getUserSecretByName(secretName)
    if err != nil {
        return "", err
    }

    decryptedSecret, err := decryptSecret(userSecret.SecretDataEncrypted, blob.getEncryptionKey(secretName, masterPassword))
    if err != nil {
        return "", errors.New("error decrypting secret: " + err.Error())
    }

    return decryptedSecret, nil
}

func (blob *UserSecretsBlob) AddOrModifySecret(secretName string, secretData string, masterPassword string, ctx context.Context) error {
    encryptedSecret, err := encryptSecret(secretData, blob.getEncryptionKey(secretName, masterPassword))
    if err != nil {
        return errors.New("error encrypting secret: " + err.Error())
    }

    userSecret, err := blob.getUserSecretByName(secretName)
    if err != nil || userSecret == nil {
        userSecret = &UserSecret{
            SecretName:secretName,
            SecretDataEncrypted:encryptedSecret,
        }
        blob.UserSecrets = append(blob.UserSecrets, userSecret)
    } else {
        userSecret.SecretDataEncrypted = encryptedSecret
    }

    // TODO: Avi: Move this somewhere else (like a set-dirty thing for transactions)
    err = storage_service.SetBlob(blob, ctx)
    if err != nil {
        logger.LogError("error saving blob after add/modify secret" +
                        "|blob name=" + kBlobName +
                        "|user id=" + strconv.FormatInt(blob.UserId, 10),
                        "|error=" + err.Error())
        return errors.New("error saving changes")
    }

    return nil
}

func (blob *UserSecretsBlob) DeleteSecret(secretName string, ctx context.Context) error {
    var index int
    for i, userSecret := range blob.UserSecrets {
        if userSecret.SecretName == secretName {
            index = i
            break
        }
    }
    if index >= len(blob.UserSecrets) {
        return errors.New("no such secret")
    }

    blob.UserSecrets = append(blob.UserSecrets[:index], blob.UserSecrets[index+1:]...)

    // TODO: Avi: Move this somewhere else (like a set-dirty thing for transactions)
    err := storage_service.SetBlob(blob, ctx)
    if err != nil {
        logger.LogError("error saving blob after deleting secret" +
                        "|blob name=" + kBlobName +
                        "|user id=" + strconv.FormatInt(blob.UserId, 10),
                        "|error=" + err.Error())
        return errors.New("error saving changes")
    }

    return nil
}

////////////////////////////////////////////////////////////////////////////////

func newUserSecretsBlob(userId int64) *UserSecretsBlob {
    userSecrets := &UserSecretsBlob{
        UserId:userId,
    }
    userSecrets.BlobDescriptor = storage_typedefs.NewBlobDescriptor(storage_typedefs.STORAGE_SPACE_APP,
                                                                    kBlobName,
                                                                    []string { "UserId" },
                                                      true)

    return userSecrets
}

func (blob *UserSecretsBlob) getUserSecretByName(secretName string) (*UserSecret, error) {
    for _, secret := range blob.UserSecrets {
        if secret.SecretName == secretName {
            return secret, nil
        }
    }

    return nil, errors.New("no such secret")
}

func (blob *UserSecretsBlob) getEncryptionKey(secretName string, masterPassword string) string {
    return secretName + ":" + strconv.FormatInt(blob.UserId, 10) + ":" + masterPassword
}

func encryptSecret(secret string, key string) (string, error) {
    encrypted, err := encryption_utils.EncryptUsingAES(secret, key)
    if err != nil {
        return "", err
    }

    return encrypted, nil
}

func decryptSecret(encryptedSecret string, key string) (string, error) {
    decrypted, err := encryption_utils.DecryptUsingAES(encryptedSecret, key)
    if err != nil {
        return "", err
    }

    return decrypted, nil
}
