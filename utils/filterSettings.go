package utils

func GetFilterSetting(voice string) string {
	switch voice {
	case "1":
		// ROBOTIC VOICE - High pitched, mechanical
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"rubberband=pitch=1.7," +
			"highpass=f=300," +
			"lowpass=f=4000," +
			"acrusher=bits=8:mode=log:aa=1," +
			"chorus=0.5:0.5:40:0.3:0.2:2," +
			"flanger=delay=3:depth=3:speed=0.5," +
			"acompressor=threshold=-15dB:ratio=6:attack=5:release=50," +
			"alimiter=limit=0.9"

	case "2":
		// DEEP DEMON VOICE - Very low, menacing
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"rubberband=pitch=0.55," +
			"equalizer=f=100:t=q:w=1:g=6," +
			"equalizer=f=200:t=q:w=1.5:g=4," +
			"lowpass=f=3500," +
			"overdrive=gain=3:colour=15," +
			"aecho=0.6:0.5:50:0.3," +
			"acompressor=threshold=-18dB:ratio=5:attack=10:release=100," +
			"alimiter=limit=0.9"

	case "3":
		// ALIEN/WHISPER - Ethereal, warbling
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"rubberband=pitch=1.35," +
			"highpass=f=400," +
			"lowpass=f=4000," +
			"chorus=0.4:0.4:45:0.3:0.2:2," +
			"aphaser=type=t:speed=0.4:decay=0.4," +
			"tremolo=f=6:d=0.4," +
			"aecho=0.6:0.5:25:0.25," +
			"acompressor=threshold=-20dB:ratio=6:attack=5:release=50," +
			"alimiter=limit=0.9"

	case "4":
		// DISTORTED ANONYMOUS - Witness protection style
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"rubberband=pitch=0.65," +
			"highpass=f=500," +
			"lowpass=f=2800," +
			"acrusher=bits=10:mode=log:aa=1," +
			"tremolo=f=20:d=0.5," +
			"chorus=0.4:0.4:30:0.3:0.2:2," +
			"acompressor=threshold=-20dB:ratio=8:attack=5:release=50," +
			"alimiter=limit=0.9"

	case "5":
		// VOCODER SYNTHETIC - Completely artificial
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"rubberband=pitch=1.3," +
			"highpass=f=350," +
			"lowpass=f=3500," +
			"acrusher=bits=6:mode=log:aa=1," +
			"vibrato=f=5:d=0.4," +
			"aphaser=type=t:speed=0.8:decay=0.4," +
			"acompressor=threshold=-18dB:ratio=8:attack=5:release=50," +
			"alimiter=limit=0.9"

	case "6":
		// MONSTER GROWL - Inhuman deep
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"rubberband=pitch=0.45," +
			"overdrive=gain=4:colour=25," +
			"equalizer=f=150:t=q:w=2:g=6," +
			"equalizer=f=400:t=q:w=1.5:g=4," +
			"lowpass=f=3000," +
			"chorus=0.5:0.4:25:0.3:0.2:1.5," +
			"aecho=0.6:0.5:60:0.35," +
			"acompressor=threshold=-15dB:ratio=6:attack=10:release=100," +
			"alimiter=limit=0.9"

	default:
		return ""
	}
}
