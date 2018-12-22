# pgnstats

parses [PGN](https://en.wikipedia.org/wiki/Portable_Game_Notation) files and extract statistics from them. handles huge files like a champ!

## usage

example: `pgnstats -f myFile.pgn -o stats.json`

help: `pgnstats -h`

## statistics

- openings tree
- heatmaps
  _ square utilization
  _ move squares
  _ checking squares
  _ capture squares
- results (white win / black win / draw / na)
- games ending with check / mate
- game length histogram
- material count histogram (using standard values)
- material difference histogram
- end game material count histogram (using standard values)
- end game material difference histogram
- castling (black/white, same/opposite)
- min/max ELO
- min/max year
- branching factor per ply
- piece tracking
