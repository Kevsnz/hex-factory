package gamedata

import ss "hextopdown/settings"

var BeltlikeParamsList = map[ss.ObjectType]BeltLikeParameters{
	ss.OBJECT_TYPE_BELT1:         {Type: ss.BELTLIKE_TYPE_NORMAL, Tier: ss.BELT_TIER_NORMAL},
	ss.OBJECT_TYPE_BELTUNDER1:    {Type: ss.BELTLIKE_TYPE_UNDER, Tier: ss.BELT_TIER_NORMAL},
	ss.OBJECT_TYPE_BELTSPLITTER1: {Type: ss.BELTLIKE_TYPE_SPLITTER, Tier: ss.BELT_TIER_NORMAL},
}

var BeltTierParamsList = [ss.BELT_TIER_COUNT]BeltTierParameters{
	ss.BELT_TIER_NORMAL:  {Speed: 18, Reach: 5},
	ss.BELT_TIER_FAST:    {Speed: 9, Reach: 7},
	ss.BELT_TIER_EXPRESS: {Speed: 6, Reach: 9},
}
