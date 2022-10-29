/**
 * Sample React Native App
 * https://github.com/facebook/react-native
 *
 * @format
 * @flow strict-local
 */

import React, {useEffect, useState} from 'react';
import {
  SafeAreaView,
  ScrollView,
  StatusBar,
  StyleSheet,
  TextInput,
  useColorScheme,
  View,
  NativeModules,
  Button,
  Text
} from 'react-native';

import {
  Colors,
  DebugInstructions,
  Header,
  LearnMoreLinks,
  ReloadInstructions,
} from 'react-native/Libraries/NewAppScreen';

const FulaModule = NativeModules.FulaModule;

const App = () => {
  const isDarkMode = useColorScheme() === 'dark';

  const backgroundStyle = {
    backgroundColor: isDarkMode ? Colors.darker : Colors.lighter,
  };

  const [res, setRes] = useState("Result")
  const [pid, setPid] = useState("QmaUMRTBMoANXqpUbfARnXkw9esfz9LP2AjXRRr7YknDAT")
  const [mas, setMas] = useState("/ip4/192.168.0.2/tcp/64658")
  const [link, setLink] = useState("bafyreihpspvqaii6nvjphmrerg4grecxcvkxaz7thfzhnu4drczbhvqpoq")

  const submit = () => {
    console.log({pid, mas, link})
    FulaModule.testGraphExchange(pid, mas, link).then(setRes).catch(console.log)
  }

  // useEffect(() => {
    
  // }, [])

  return (
    <SafeAreaView style={backgroundStyle}>
      <StatusBar
        barStyle={isDarkMode ? 'light-content' : 'dark-content'}
        backgroundColor={backgroundStyle.backgroundColor}
      />
      <ScrollView
        contentInsetAdjustmentBehavior="automatic"
        style={backgroundStyle}>
        <Header />
        <View
          style={{
            backgroundColor: isDarkMode ? Colors.black : Colors.white,
          }}>
          <TextInput value={pid} onChangeText={setPid}  placeholder="Peer ID" />
          <TextInput value={mas} onChangeText={setMas}  placeholder="Multi Addr"/>
          <TextInput value={link} onChangeText={setLink} placeholder="Link" />
          <Button onPress={submit} title="Fetch" />
          <Text>{res}</Text>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  sectionContainer: {
    marginTop: 32,
    paddingHorizontal: 24,
  },
  sectionTitle: {
    fontSize: 24,
    fontWeight: '600',
  },
  sectionDescription: {
    marginTop: 8,
    fontSize: 18,
    fontWeight: '400',
  },
  highlight: {
    fontWeight: '700',
  },
  bordered: {
    border: '1px solid red'
  }
});

export default App;
