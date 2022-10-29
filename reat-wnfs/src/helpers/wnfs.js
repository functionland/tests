export async function init(wn) {

const state = await wn.initialise({
  permissions: {
    // Will ask the user permission to store
    // your apps data in `private/Apps/Nullsoft/Winamp`
    app: {
      name: "Winamp",
      creator: "Nullsoft"
    },

    // Ask the user permission to additional filesystem paths
    fs: {
      private: [ wn.path.directory("Party", "pics") ],
      public: [ wn.path.directory("Docs") ]
    }
  }

}).catch(err => {
    console.log({err})
  switch (err) {
    case wn.InitialisationError.InsecureContext:
      // We need a secure context to do cryptography
      // Usually this means we need HTTPS or localhost

    case wn.InitialisationError.UnsupportedBrowser:
      // Browser not supported.
      // Example: Firefox private mode can't use indexedDB.
  }

})

console.log({ssss: state})


switch (state.scenario) {

  case wn.Scenario.AuthCancelled:
    // User was redirected to lobby,
    // but cancelled the authorisation
    break;

  case wn.Scenario.AuthSucceeded:
  case wn.Scenario.Continuation:
    // State:
    // state.authenticated    -  Will always be `true` in these scenarios
    // state.newUser          -  If the user is new to Fission
    // state.throughLobby     -  If the user authenticated through the lobby, or just came back.
    // state.username         -  The user's username.
    //
    // ☞ We can now interact with our file system (more on that later)
    return state
    break;

  case wn.Scenario.NotAuthorised:
    wn.redirectToLobby(state.permissions)
    break;

}
}