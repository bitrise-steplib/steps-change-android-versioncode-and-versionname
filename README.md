# Change Android versionCode and versionName

[![Step changelog](https://shields.io/github/v/release/bitrise-community/steps-change-android-versioncode-and-versionname?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-community/steps-change-android-versioncode-and-versionname/releases)

Updates the Android versionCode and versionName in your project's `build.gradle` file.

<details>
<summary>Description</summary>

Modifies the version information of your Android app by updating versionCode and versionName attributes in your project's `build.gradle` file before you'd publish your app to Google Play Store.

### Configuring the Step

1. Insert the **Change Android versionCode and versionName** Step before a build Step such as **Android Build** or **Gradle Runner** in your Workflow.
2. Click the Step to modify its input fields.
3. Add the file path to the **Path to the build.gradle file** so that the Step knows where to find the file that contains the versionCode and versionName attributes.
4. Provide a new versionName in the **New versionName** input. If you leave this input empty, the previous versionName will be displayed on Google Play Store. The input's value must be a string in this format `<major>.<minor>.<point>`.
5. Provide a versionCode in the **New versionCode** input to track app versions. If you leave this input empty, you will release your app with the version that is already set in the `build.gradle` file's `versionCode` attribute. The input's value must be an integer.
If you wish to offset the value you set in the **New versionCode** input, then you can use the **versionCode Offset** input to add the offset integer value here. This number will be added to the versionCode's value.

### Troubleshooting

The **Change Android versionCode and versionName** Step must be inserted BEFORE the **Android Build** Step as the former makes sure you  upload the build with the right versionCode and versionName to Google Play Store.

### Useful links

- [About Android versionCode and versionNumber](https://developer.android.com/studio/publish/versioning)
- [Build versioning](https://devcenter.bitrise.io/builds/build-numbering-and-app-versioning/)

### Related Steps

- [Android Sign](https://www.bitrise.io/integrations/steps/sign-apk)
- [Android Build](https://www.bitrise.io/integrations/steps/android-build)
- [Gradle Runner](https://www.bitrise.io/integrations/steps/gradle-runner)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `build_gradle_path` | Path to the build.gradle file shows the versionCode and versionName settings. | required | `$BITRISE_SOURCE_DIR/app/build.gradle` |
| `new_version_name` | New versionName to set. Specify a string value, such as `"1.0.0"`. If the specified value is not surranded by double quote (`"`) characters, the step will add them. Leave this input empty so that versionName remains unchanged. |  |  |
| `new_version_code` | New versionCode to set. Specify a positive integer value, such as `1`. The greatest value Google Play allows for versionCode is 2100000000. Clear this input's default value to leave the versionCode unchanged. |  | `$BITRISE_BUILD_NUMBER` |
| `version_code_offset` | Offset value to add to `New versionCode`, for example: `1`. Leave this input empty if you want the exact value you set in `New versionCode` input. |  |  |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `ANDROID_VERSION_NAME` | The final Android versionName set in the build.gradle file after the Step execution. |
| `ANDROID_VERSION_CODE` | The final Android versionCode set in the build.gradle file after the Step execution. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-community/steps-change-android-versioncode-and-versionname/pulls) and [issues](https://github.com/bitrise-community/steps-change-android-versioncode-and-versionname/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
