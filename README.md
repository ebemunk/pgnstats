pgnstats
--------

parses [PGN](https://en.wikipedia.org/wiki/Portable_Game_Notation) files and extract statistics from them. handles huge files like a champ! mostly a companion to [chess-dataviz](https://github.com/ebemunk/chess-dataviz)

### usage
`./pgnstats -h` for help
`./pgnstats -f=myFile.pgn -o=stats.json`

### statistics
* openings tree
* heatmaps
	* square utilization
	* move squares
	* checking squares
	* capture squares
* results (white win / black win / draw / na)
* games ending with check / mate
* game length histogram
* material count histogram (using standard values)
* material difference histogram
* castling (black/white, same/opposite)
* min/max ELO
* min/max year
