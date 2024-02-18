<script>
  import { selectedAlbumStore } from '../stores/selectedAlbumStore';
  import albumArt from 'album-art';

  let artworkUrl;
  let tracksDisplay = [];

  selectedAlbumStore.subscribe(async album => {
    if (album) {
      artworkUrl = await albumArt(album?.Artist, {album: album?.Album});
      setTracksDisplay(album);
    }
  });

  function setTracksDisplay(album) {
    tracksDisplay = [];
    const tracklist = JSON.parse(album.Tracklist);
    for (let track of tracklist) {
      tracksDisplay.push({num: track.Number, title: track.Title});
    }
  }
</script>

{#if artworkUrl}
  <img src={artworkUrl} alt="album cover" />
  <div style="text-align: left;">
    {#each tracksDisplay as track}
      <p>{track?.num} - {track?.title}</p>
    {/each}
  </div>
{/if}

<style>
  img {
    width: 75%;
    margin: 20px auto;
  }
</style>
