# Setting up SOPS encryption

Follow these instructions to set up encryption of secrets in this repo using SOPS, and set up decryption of those secrets in your cluster.

For full details, see [this docs page on the Flux site](https://fluxcd.io/docs/guides/mozilla-sops/)

## Encrypting files in the repo

1. Install `sops` using the instructions [here](https://fluxcd.io/docs/guides/mozilla-sops/#prerequisites)

2. Generate a GPG key using the instructions [here](https://fluxcd.io/docs/guides/mozilla-sops/#generate-a-gpg-key). You can skip this step if you already have a GPG key you want to use.

3. Export your Key Fingerprint as an environment variable. We'll use it for a few of the subsequent steps.

    ```shell
    export KEY_FP="YourKeyFingerprintHere"
    ```

4. Export the public key into the Git repository so your teammates can use it to encrypt files. Commit and push the file to your repo.

    ```shell
    # Create the public key file
    gpg --export --armor "${KEY_FP}" > ./.sops.pub.asc
    # Teammates run this to import the key
    gpg --import ./.sops.pub.asc
    ```

5. Create the SOPS configuration file. Commit and push the file to your repo.

    ```shell
    cat <<EOF > ./.sops.yaml
    creation_rules:
      - pgp: "${KEY_FP}"
    EOF
    ```

6. Encrypt/Decrypt files (Note: to decrypt you need the private key too)

```shell
# Encrypt the file
sops -e -i thefile.yaml
# Decrypt the file
sops -d -i thefile.yaml
```

## Configure the cluster to decrypt

1. Create a secret in the `flux-system` namespace that contains the GPG secret key. This needs to be done AFTER you've deployed Flux.

    ```shell
    gpg --export-secret-keys --armor "${KEY_FP}" | kubectl create secret generic sops-gpg --namespace=flux-system --from-file=sops.asc=/dev/stdin
    ```

    If you need to sneakernet the secret do it this way:

    ```shell
    # Create the file
    gpg --export-secret-keys --armor "${KEY_FP}" | kubectl create secret generic sops-gpg --dry-run=client -o=yaml --namespace=flux-system --from-file=sops.asc=/dev/stdin > secret-sops-gpg.yaml

    # Sneakernet it across

    # Deploy it
    kubectl apply -f secret-sops-gpg.yaml
    ```
