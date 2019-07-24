import axios from 'axios'
import qs from 'qs'
import Cookies from 'js-cookie'

function gettoken () {
  return Cookies.get('appdaemon.token')
}

function storetoken (token) {
  Cookies.set('appdaemon.token', token, { path: '/', expires: 1 })
}

function removetoken () {
  Cookies.remove('appdaemon.token')
}

function geturl (name, params) {
  let baseurl = process.env.VUE_APP_API_BASE_URL
  if (baseurl) {
    baseurl = baseurl.replace(/\/+$/, '')
  } else {
    baseurl = 'http://localhost:7982'
  }
  return baseurl + '/' + name + '?' + qs.stringify(params)
}

async function login (username, password) {
  let r = await axios.post(geturl('Login'), qs.stringify({ username, password }))
  if (r.data.code === 1) {
    storetoken(r.data.data.token)
  }
  return r.data
}

async function getloginuser () {
  let r = await axios.get(geturl('GetLoginUser', { token: gettoken() }))
  return r.data
}

async function getluascriptlist (page, pagesize, keyword, addr) {
  let r = await axios.get(geturl('GetLuaScriptList', { token: gettoken(), page, pagesize, keyword, addr }))
  return r.data
}

async function getluascript (id) {
  let r = await axios.get(geturl('GetLuaScript', { token: gettoken(), id }))
  return r.data
}

async function addluascript (name, script, status, addr) {
  let r = await axios.post(geturl('AddLuaScript', { token: gettoken() }), qs.stringify({ name, script, status, addr }))
  return r.data
}

async function updateluascript (id, name, script, status, addr) {
  let r = await axios.post(geturl('UpdateLuaScript', { token: gettoken() }), qs.stringify({ id, name, script, status, addr }))
  return r.data
}

async function updateluascriptstatus (id, status) {
  let r = await axios.post(geturl('UpdateLuaScriptStatus', { token: gettoken() }), qs.stringify({ id, status }))
  return r.data
}

async function deleteluascript (id) {
  let r = await axios.post(geturl('DeleteLuaScript', { token: gettoken() }), qs.stringify({ id }))
  return r.data
}

async function getserverlist () {
  let r = await axios.get(geturl('GetServerList', { token: gettoken() }))
  return r.data
}

export default {
  gettoken,
  storetoken,
  removetoken,
  geturl,
  login,
  getloginuser,
  getluascriptlist,
  getluascript,
  addluascript,
  updateluascript,
  updateluascriptstatus,
  deleteluascript,
  getserverlist
}
