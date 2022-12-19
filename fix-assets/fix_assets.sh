#! /bin/bash

set -e

if [ -z "${ASSET_PROFILE_PATH}" ]; then
	echo '$ASSET_PROFILE_PATH not set'
	exit 1
fi

# if [ -z "${ICON_GEN_CONFIG_PATH}" ]; then
	# echo '$ICON_GEN_CONFIG_PATH not set'
	# exit 1
# fi

rings_list=$(curl -s https://rings-api.pando.im/api/v1/markets)
for r in $(echo "${rings_list}" | jq -r '.[] | @base64'); do
		_jq() {
		 echo ${r} | base64 --decode | jq -r ${1}
		}
		 ctoken_asset_id=$(_jq '.ctoken_asset_id')
		 asset_id=$(_jq '.asset_id')
		if [ ! -d "$ASSET_PROFILE_PATH/assets/$ctoken_asset_id" ]; then
			echo $ctoken_asset_id $asset_id
			$(icon-gen rings --a $ctoken_asset_id --a1 $asset_id --config $ICON_GEN_CONFIG_PATH)
		fi
done

fswap_list=$(curl -s https://api.4swap.org/api/pairs)
for r in $(echo "${fswap_list}" | jq -r '.data.pairs[] | @base64'); do
		_jq() {
		 echo ${r} | base64 --decode | jq -r ${1}
		}
		 liquidity_asset_id=$(_jq '.liquidity_asset_id')
		 base_asset_id=$(_jq '.base_asset_id')
		 quote_asset_id=$(_jq '.quote_asset_id')
		if [ ! -d "$ASSET_PROFILE_PATH/assets/$liquidity_asset_id" ]; then
			$(icon-gen fswap --a $liquidity_asset_id --a1 $base_asset_id --a2 $quote_asset_id --config $ICON_GEN_CONFIG_PATH)
		fi
done

cp -r ./outputs/* $ASSET_PROFILE_PATH/assets
