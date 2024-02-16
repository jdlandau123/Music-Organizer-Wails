<script>
  import logo from './assets/images/logo-universal.png';
  import Header from './components/Header.svelte';
  import CollectionTable from './components/CollectionTable.svelte';
  import { SyncMusicCollection } from '../wailsjs/go/main/App.js';

  let err = '';
  async function syncMusicCollection() {
    try {
      const err = await SyncMusicCollection();
    } catch(e) {
      err = e;
    }
  }

</script>

<main>
  <Header />
  <CollectionTable />
  <button on:click={syncMusicCollection}>Sync Music Collection</button>
  <p>{err}</p>
</main>

<style>
  main {
    max-width: 100vw;
    max-height: 100vh;
    overflow: hidden
  }
</style>
