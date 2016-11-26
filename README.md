These are tools that I found useful for myself when working with Amazon Web Services.

### How to use it
* In order to use it you need to have `~/.aws/credentials` file with AWS credentials on your computer (Linux/Mac). Futher information can be found [here](https://aws.amazon.com/blogs/security/a-new-and-standardized-way-to-manage-credentials-in-the-aws-sdks/).
* Download aws-cli-tools binary in the root.

### SSH access for instances in Opsworks
  `./aws-cli-tools -command ssh [-stack name] [-layer name] [-prefix prefix]`

stack, layer and prefix are all optional.
stack - stack name in lower case, with underline symbols instead of spaces.
layer - layer name.
prefix - stack name prefix. Useful if all your stacks have the same prefix, you just specify prefix and stack name without it.

If you do not specify any params - you will see a dialog with available stack names and layers to choose from.

For example you have a stacks with "Delta Gamma" and "Delta Beta" names. You could add an alias `delta` for `./aws-cli-tools -command ssh -prefix delta` and then access both stacks with `delta -stack gamma` and `delta -stack beta` commands.

### RDS instances information
  `./aws-cli-tools --command dbinfo [-rds rds_name]`

rds - is an optional parameter, that will show info for a particular rds instance. Without it information for all rds instances will be given.

Will display info for each rds instance. It's status, instance type, freeable memory, free storage space. As long as cpu utilization, write iops, read iops for the past 10 minutes with 1 minute interval.
