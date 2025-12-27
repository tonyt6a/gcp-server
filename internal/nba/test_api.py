from nba_api.live.nba.endpoints import scoreboard
from nba_api.stats.endpoints import playbyplay, playbyplayv2, playbyplayv3

# # Today's Score Board
# games = scoreboard.ScoreBoard()

# # json
# games.get_json()

# # dictionary
# d = games.get_dict()

# print(d)

'''
v3 probably more helpful, punt on v1
'''
# pbp
# pbp1 = playbyplay.PlayByPlay(game_id="0021700807")
# pbp1.get_json()
# d1 = pbp1.get_dict()
# # resultSets has pbp
# rs = d1['resultSets']
# # Only care about index 0
# test1 = rs[0]
# # Only care about rowSet
# rowset = test1['rowSet']
# for r in rowset:
#     print(r)
# try:
#     pbp2 = playbyplayv2.PlayByPlayV2(game_id="0021700807")
#     pbp2.get_json()
#     d2 = pbp2.get_dict()
#     print(d2)
# except KeyError:
#     print("Key Error")

import pandas as pd
import matplotlib as plt

def time_remaining(time:str):
    # Time will be of format PTXXMYY.YYS
    # Unsure of what PT means but will ignore for now
    # Returns number of time left in quarter in seconds
    minute_range = (2, 4)
    seconds_range = (5, 10)
    minutes = int  (time[minute_range[0] : minute_range[1]])
    seconds = float(time[seconds_range[0]: seconds_range[1]])
    return (minutes * 60) + seconds


# pbp 3
pbp3 = playbyplayv3.PlayByPlayV3(game_id="0021700807")
pbp3.get_json()
d3 = pbp3.get_dict()
# game has pbp
game = d3['game']
# Only care about actions
actions = game['actions']
d = {}
for a in actions:
    for k, v in a.items():
        if k in d:
            d[k].append(v)
        else:
            d[k] = [v]
df = pd.DataFrame(d)

print(df)
