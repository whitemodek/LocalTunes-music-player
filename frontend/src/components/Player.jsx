import { useState, useRef } from 'react';
import { Box, IconButton, Typography } from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import PauseIcon from '@mui/icons-material/Pause';
import ReactPlayer from 'react-player';

export default function Player({ track }) {
  const [playing, setPlaying] = useState(true);
  const audioRef = useRef(null);

  const togglePlay = () => {
    if (playing) {
      audioRef.current?.pause();
    } else {
      audioRef.current?.play();
    }
    setPlaying(!playing);
  };

  return (
    <Box sx={{
      position: 'fixed', bottom: 0, left: 0, right: 0, bgcolor: '#282828',
      display: 'flex', alignItems: 'center', p: 2, borderTop: '1px solid #444', zIndex: 1000
    }}>
      <ReactPlayer 
        ref={audioRef} 
        url={`http://localhost:8080${track.streamUrl}`} 
        playing={playing} 
        hidden 
        width={0} height={0}
      />
      
      <Box sx={{ display: 'flex', alignItems: 'center', width: '30%' }}>
        <img src={track.coverUrl || 'https://via.placeholder.com/50'} alt="" style={{ width: 50, height: 50, marginRight: 10 }} />
        <Box>
          <Typography fontWeight="bold">{track.title}</Typography>
          <Typography variant="body2" color="grey.500">{track.artist}</Typography>
        </Box>
      </Box>

      <Box sx={{ display: 'flex', justifyContent: 'center', width: '40%' }}>
        <IconButton onClick={togglePlay}>
          {playing ? <PauseIcon sx={{ color: 'white' }} /> : <PlayArrowIcon sx={{ color: 'white' }} />}
        </IconButton>
      </Box>
    </Box>
  );
}