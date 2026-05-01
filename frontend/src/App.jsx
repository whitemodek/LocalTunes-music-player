import { useState, useEffect } from 'react';
import axios from 'axios';
import { Box, TextField, List, ListItem, ListItemAvatar, Avatar, ListItemText, Typography, Button, Container } from '@mui/material';
import Player from './components/Player';

function App() {
  const [tracks, setTracks] = useState([]);
  const [query, setQuery] = useState('');
  const [currentTrack, setCurrentTrack] = useState(null);

  const fetchTracks = async (q = '') => {
    try {
      const res = await axios.get(`http://localhost:8080/api/search?q=${q}`);
      setTracks(res.data);
    } catch (err) {
      console.error('Ошибка загрузки треков:', err);
    }
  };

  useEffect(() => { fetchTracks(); }, []);

  const handleUpload = async (e) => {
    const formData = new FormData();
    for (let f of e.target.files) formData.append('track', f);
    
    try {
      await axios.post('http://localhost:8080/api/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      fetchTracks(query);
      e.target.value = ''; // Сбрасываем input
    } catch (err) {
      console.error('Ошибка загрузки:', err);
    }
  };

  return (
    <Box sx={{ bgcolor: '#1a1a1a', color: 'white', minHeight: '100vh', pb: 12 }}>
      <Container maxWidth="lg">
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, py: 2, borderBottom: '1px solid #333' }}>
          <Typography variant="h4" sx={{ color: '#FFCC00', flex: 1 }}>LocalTunes</Typography>
          <TextField
            placeholder="🔍 Поиск по артисту или треку..."
            value={query}
            onChange={(e) => { setQuery(e.target.value); fetchTracks(e.target.value); }}
            sx={{ input: { color: 'white' }, bgcolor: '#333', borderRadius: 1, p: 1 }}
          />
          <Button variant="contained" component="label" sx={{ bgcolor: '#FFCC00', color: 'black' }}>
            Загрузить
            <input type="file" hidden multiple accept="audio/*" onChange={handleUpload} />
          </Button>
        </Box>

        <List>
          {tracks.map(track => (
            <ListItem 
              key={track.id} 
              onClick={() => setCurrentTrack(track)}
              sx={{ borderBottom: '1px solid #333', '&:hover': { bgcolor: '#282828' }, cursor: 'pointer' }}
            >
              <ListItemAvatar>
                <Avatar src={track.coverUrl || 'https://via.placeholder.com/60'} variant="square" sx={{ width: 60, height: 60, mr: 2 }} />
              </ListItemAvatar>
              <ListItemText 
                primary={track.title} 
                secondary={track.artist}
                primaryTypographyProps={{ color: 'white' }}
                secondaryTypographyProps={{ color: 'grey.500' }}
              />
            </ListItem>
          ))}
        </List>
      </Container>

      {currentTrack && <Player track={currentTrack} />}
    </Box>
  );
}
export default App;
