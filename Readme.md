# go-resimg

指定したディレクトリ以下のファイル構造を維持しつつ、ターゲットディレクトリに画像をリサイズして保存する。

画像はアスペクト比を保ちながら正方形のjpg画像として出力される。余った余白は黒塗りになる。

Usage: go-resimg [Option] dir target  
 -p int  
	-p [process-num] (default 12)  
 -s int  
	-s [image-size] (default 100)