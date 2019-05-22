(debug)
(load "../lib/all.gf")

(env fire (width 50 height 25
           buf (new-bin (* width height 1))
           esc (str 0x1b "[")
           out stdout
           max-fade 50
           avg-n 3 avg-time _)
  (fun get-offs (x y)
    (+ (- width x 1) (* (- height y 1) width)))

  (fun xy (x y)
    (# buf (get-offs x y)))

  (fun set-xy (f x y)
    (let i (get-offs x y))
    (set (# buf i) (f (# buf i)))
    _)

  (fun clear ()
    (print out (str esc "2J")))

  (fun home ()
    (print out (str esc "H")))

  (fun pick-color (r g b)
    (print out (str esc "48;2;" (int r) ";" (int g) ";" (int b) "m")))

  (fun restore-color ()
    (print out (str esc "0m")))

  (fun restore ()
    (clear)
    (home)
    (flush out))

  (fun init ()
    (for (width x)
      (set (xy x 0) 0xff))

    (clear))

  (fun render ()
    (let t0 (now))
    
    (for ((- height 1) y)
      (for (width x)
        (let v (xy x y))
        (if (and x (< x (- width 1))) (inc x (- 1 (rand 3))))
        (set (xy x (+ y 1)) (- v (rand (min max-fade (+ (int v) 1)))))))

    (let i -1)
    (home)
    
    (for (height y)
      (for (width x)
        (let g (# buf (inc i)) r (if g 0xff 0) b (if (= g 0xff) 0xff 0))
        (pick-color r g b)
        (print out " "))

      (print out \n))

    (set avg-time (if (_? avg-time)
                    (- (now) t0)
                    (+ (/ (* avg-time (- avg-n 1)) avg-n)
                       (/ (- (now) t0) avg-n))))
    
    (restore-color)
    (print out (/ 1000000000.0 avg-time))
    (flush out)))

(fire/init)

(for 50
  (fire/render))

(fire/restore)