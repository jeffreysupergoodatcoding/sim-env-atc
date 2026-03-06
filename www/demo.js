// Demo Service Frontend
class DemoService {
    constructor() {
        this.clips = [];
        this.currentIndex = -1;
        this.isPlaying = false;
        this.audio = new Audio();
        this.demoMode = false;

        this.init();
    }

    async init() {
        try {
            // Check if demo mode is enabled in config
            const configResponse = await fetch('/api/v1/config');
            if (configResponse.ok) {
                const config = await configResponse.json();
                this.demoMode = config.demo_mode === true;
                console.log('[Demo] Demo mode status:', this.demoMode);
            }

            if (!this.demoMode) return;

            // Fetch demo clips
            const response = await fetch('/api/v1/demo/clips');
            if (response.ok) {
                this.clips = await response.json();
                console.log('[Demo] Loaded clips:', this.clips);
            }

            // Listen for audio ended to update state
            this.audio.onended = () => {
                this.isPlaying = false;
                this.triggerReactivity();
            };

        } catch (error) {
            console.error('[Demo] Initialization failed:', error);
        }
    }

    async nextClip() {
        if (!this.demoMode) return;

        try {
            const response = await fetch('/api/v1/demo/next', { method: 'POST' });
            if (response.ok) {
                const clip = await response.json();
                this.currentIndex = this.clips.findIndex(c => c.audio_url === clip.audio_url);

                // Play audio
                this.playAudio(clip.audio_url);

                // Trigger AI Advisory visualization
                this.simulateAIAdvisory(clip);

                this.triggerReactivity();
                return clip;
            }
        } catch (error) {
            console.error('[Demo] Failed to get next clip:', error);
        }
    }

    async reset() {
        if (!this.demoMode) return;

        try {
            await fetch('/api/v1/demo/reset', { method: 'POST' });
            this.currentIndex = -1;
            this.isPlaying = false;
            this.audio.pause();

            // Reset pipeline state
            const store = window.Alpine.store('atc');
            if (store) {
                store.demoPipeline = {
                    transcription: '',
                    extraction: null,
                    verification: null,
                    response: '',
                    status: 'idle'
                };
            }

            this.triggerReactivity();

            // Refresh aircraft list
            if (window.Alpine && window.Alpine.store('atc')) {
                window.Alpine.store('atc').fetchAircraft();
            }
        } catch (error) {
            console.error('[Demo] Failed to reset demo:', error);
        }
    }

    playAudio(url) {
        this.audio.src = url;
        this.audio.play();
        this.isPlaying = true;
        this.triggerReactivity();
    }

    async simulateAIAdvisory(clip) {
        // This simulates the NLU/NLG pipeline in the UI with realistic steps
        console.log('[Demo] Simulating AI Advisory for:', clip.transcription);

        // Reset store state
        const store = window.Alpine.store('atc');
        store.demoPipeline = {
            transcription: '',
            extraction: null,
            verification: null,
            response: '',
            status: 'active'
        };

        // Step 1: Transcription (Whisper) - Type out the words
        const words = clip.transcription.split(' ');
        let currentText = '';
        for (let i = 0; i < words.length; i++) {
            currentText += words[i] + ' ';
            store.demoPipeline.transcription = currentText;
            await new Promise(r => setTimeout(r, 80 + Math.random() * 120)); // Faster typing for real demo
        }

        // Step 2: NLU Extraction (Mistral 7B)
        await new Promise(r => setTimeout(r, 600));
        store.demoPipeline.status = 'extraction';
        store.demoPipeline.extraction = {
            callsign: clip.entities.callsign || '-',
            intent: clip.entities.instruction_type || 'INSTRUCTION',
            facility: clip.entities.facility || '-',
            confidence: clip.entities.confidence || '0.95'
        };

        // Step 3: Safety Gate (ADS-B Check)
        await new Promise(r => setTimeout(r, 800));
        store.demoPipeline.status = 'verification';

        // Check if callsign exists in simulated aircraft
        // The store.aircraft is a list of aircraft objects, not a map by callsign
        // We need to look it up correctly
        let matchingAircraft = null;
        if (store.aircraft) {
            const cleanCallsign = clip.callsign.replace(/\s+/g, '').toUpperCase();
            matchingAircraft = store.aircraft.find(a =>
                a.flight && a.flight.replace(/\s+/g, '').toUpperCase().includes(cleanCallsign)
            );
        }

        if (clip.verified && matchingAircraft) {
            store.demoPipeline.verification = {
                status: 'verified',
                details: `MATCH: ${clip.callsign} at ${matchingAircraft.lat.toFixed(4)}, ${matchingAircraft.lon.toFixed(4)}`
            };
        } else {
            store.demoPipeline.verification = {
                status: 'failed',
                details: `ERROR: UNABLE TO VERIFY ${clip.callsign}`
            };
        }

        // Step 4: NLG Readback
        if (store.demoPipeline.verification.status === 'verified') {
            await new Promise(r => setTimeout(r, 1000));
            store.demoPipeline.status = 'response';
            store.demoPipeline.response = clip.response || `Roger, ${clip.callsign}, instruction acknowledged.`;
        } else {
            await new Promise(r => setTimeout(r, 400));
            store.demoPipeline.status = 'response';
            store.demoPipeline.response = clip.response || "Frequency busy. Say again.";
        }

        // Step 5: Log to RLHF (Simulated)
        const logEntry = {
            id: Date.now(),
            timestamp: new Date().toISOString(),
            clip_index: clip.index,
            transcription: clip.transcription,
            extraction: store.demoPipeline.extraction,
            verification: store.demoPipeline.verification,
            response: store.demoPipeline.response,
            reward: null // Placeholder for human scoring
        };

        if (!store.pipelineLogs) store.pipelineLogs = [];
        store.pipelineLogs.unshift(logEntry);

        // Dispatch event for UI updates
        document.dispatchEvent(new CustomEvent('pipeline-log-added', { detail: logEntry }));
    }

    triggerReactivity() {
        const store = window.Alpine.store('atc');
        if (store) {
            store.demoCurrentIndex = this.currentIndex;
            store.demoClips = this.clips;
            store.demoIsPlaying = this.isPlaying;
        }
        document.dispatchEvent(new CustomEvent('demo-state-updated'));
    }
}

// Initialize and export
window.demoService = new DemoService();
