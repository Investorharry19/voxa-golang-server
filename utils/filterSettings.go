package utils

func GetFilterSetting(voice string) string {
	switch voice {
	case "1":
		// High-pitched robotic voice
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"highpass=f=80," +
			"lowpass=f=12000," +
			"acompressor=threshold=-20dB:ratio=4:attack=5:release=50," +
			"rubberband=pitch=1.3348," +
			"highpass=f=200," +
			"equalizer=f=2500:t=q:w=1.5:g=5," +
			"equalizer=f=6000:t=q:w=1:g=3," +
			"chorus=0.6:0.9:50:0.4:0.25:2"

	case "2":
		// Thin, chipmunk-like voice
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"highpass=f=80," +
			"lowpass=f=12000," +
			"acompressor=threshold=-18dB:ratio=3:attack=5:release=50," +
			"rubberband=pitch=1.2599," +
			"highpass=f=250," +
			"equalizer=f=3000:t=q:w=1.5:g=6," +
			"equalizer=f=7000:t=q:w=1:g=4," +
			"equalizer=f=9000:t=q:w=2:g=-3"

	case "3":
		// Deep, mysterious voice
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"highpass=f=60," +
			"lowpass=f=12000," +
			"acompressor=threshold=-22dB:ratio=6:attack=5:release=100," +
			"rubberband=pitch=0.7492," +
			"equalizer=f=100:t=q:w=1:g=4," +
			"equalizer=f=300:t=q:w=1.5:g=5," +
			"equalizer=f=1800:t=q:w=2:g=-5," +
			"lowpass=f=4500," +
			"aecho=0.8:0.9:40:0.3"

	case "4":
		// Muffled, distant voice
		return "loudnorm=I=-16:TP=-1.5:LRA=11," +
			"highpass=f=80," +
			"acompressor=threshold=-20dB:ratio=4:attack=5:release=50," +
			"rubberband=pitch=0.8909," +
			"highpass=f=400," +
			"lowpass=f=3400," +
			"equalizer=f=1200:t=q:w=1.5:g=4," +
			"overdrive=gain=2:colour=5," +
			"alimiter=limit=0.8:attack=5:release=50"

	default:
		return ""
	}
}
