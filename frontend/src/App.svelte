<script>
  import logo from './assets/images/logo-universal.png';
  import Header from './components/Header.svelte';
  import CollectionTable from './components/CollectionTable.svelte';
  import Snackbar from './components/Snackbar.svelte';
  import { SyncMusicCollection } from '../wailsjs/go/main/App.js';

  let snackbar;

  async function syncMusicCollection() {
    try {
      const err = await SyncMusicCollection();
    } catch(e) {
      snackbar.open(e);
    }
  }

  function syncDevice() {
    console.log('syncing');
  }

</script>

<main>
  <Header />
  <div style="display: flex;">
    <CollectionTable />
    <div style="flex: 1; padding: 20px;">
      <button class="secondaryBtn" on:click={syncMusicCollection}>
        Sync Music Collection
      </button>
      <button class="primaryBtn" on:click={syncMusicCollection}>
        Sync Device
      </button>
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
</style>
