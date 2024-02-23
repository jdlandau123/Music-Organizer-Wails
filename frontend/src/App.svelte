<script>
  import Header from './components/Header.svelte';
  import CollectionTable from './components/CollectionTable.svelte';
  import Snackbar from './components/Snackbar.svelte';
  import { SyncMusicCollection, ScanDevice, SyncDevice } from '../wailsjs/go/main/App.js';
  import AlbumDetails from './components/AlbumDetails.svelte';
  import Loader from './components/Loader.svelte';
  import { albumsStore } from './stores/albumsStore';

  let snackbar;
  let collectionTable;
  let showLoading = false;

  async function syncMusicCollection() {
    showLoading = true;
    try {
      await SyncMusicCollection();
      showLoading = false;
    } catch(e) {
      snackbar.open(e);
    }
    collectionTable.getAlbums();
  }

  async function scanDevice() {
    showLoading = true;
    try {
      await ScanDevice();
      showLoading = false;
    } catch(e) {
      snackbar.open(e);
    }
    collectionTable.getAlbums();
  }

  async function syncDevice() {
    showLoading = true;
    const ids = $albumsStore.filter(a => a.IsOnDevice === true).map(a => a.Id);
    try {
      await SyncDevice(ids);
      showLoading = false;
    } catch(e) {
      snackbar.open(e);
    }
  }

</script>

<main>
  <Header />
  <div class="wrapper">
    {#if showLoading}
      <Loader />
    {:else}  
      <CollectionTable bind:this={collectionTable} />
    {/if}
    <div class="controls">
      <button class="secondaryBtn tooltip" on:click={syncMusicCollection}>
        Sync Music Collection
        <span class="tooltiptext">
          Sync the app database with your music collection
        </span>
      </button>
      <button class="secondaryBtn tooltip" on:click={scanDevice}>
        Scan Device
        <span class="tooltiptext">
          Sync the app database with the music you currently have on your device
        </span>
      </button>
      <button class="primaryBtn tooltip" on:click={syncDevice}>
        Sync Device
        <span class="tooltiptext">
          Transfer the selected albums to your device
        </span>
      </button>
      <AlbumDetails />
    </div>
  </div>
  <Snackbar bind:this={snackbar} />
</main>

<style>
  main {
    max-width: 100vw;
    max-height: 100vh;
    overflow: hidden;
  }

  button {
    margin-left: 10px;
    margin-right: 10px;
  }

  .wrapper {
    display: flex;
    justify-content: space-between;
    gap: 20px;
    padding: 20px;
  }

  .controls {
    height: calc(100vh - 120px);
    overflow-y: auto;
    overflow-x: hidden;
    width: 40%;
  }

  .tooltip {
    position: relative;
    display: inline-block;
  }

  .tooltip .tooltiptext {
    visibility: hidden;
    width: 150px;
    background-color: rgba(0, 0, 0, 0.8);
    color: #fff;
    text-align: center;
    border-radius: 8px;
    padding: 5px 0;
    position: absolute;
    z-index: 1;
    top: 100%;
    left: 50%;
    margin-left: -60px;
    margin-top: 10px;
  }

  .tooltip:hover .tooltiptext {
    visibility: visible;
  }
</style>
