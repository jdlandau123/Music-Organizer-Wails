<script>
  import Header from './components/Header.svelte';
  import CollectionTable from './components/CollectionTable.svelte';
  import Snackbar from './components/Snackbar.svelte';
  import { SyncMusicCollection, SyncDevice } from '../wailsjs/go/main/App.js';
  import AlbumDetails from './components/AlbumDetails.svelte';

  let snackbar;

  async function syncMusicCollection() {
    try {
      await SyncMusicCollection();
    } catch(e) {
      snackbar.open(e);
    }
  }

  async function syncDevice() {
    try {
      await SyncDevice();
    } catch(e) {
      snackbar.open(e);
    }
  }

</script>

<main>
  <Header />
  <div class="wrapper">
    <CollectionTable />
    <div style="height: calc(100vh - 90px); overflow-y: auto; overflow-x: hidden;">
      <button class="secondaryBtn" on:click={syncMusicCollection}>
        Sync Music Collection
      </button>
      <button class="primaryBtn" on:click={syncDevice}>
        Sync Device
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
    overflow: hidden
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
</style>
