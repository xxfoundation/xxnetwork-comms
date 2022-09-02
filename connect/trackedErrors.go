////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package connect

// This file contains errors that are tracked by or returned by the host
// depending on their tracked metrics.

// ProxyError is part of the error reported by gateway when host cannot be
// reached. Its frequency is being tracked so that it can be excluded from the
// host pool in the layer after a set number of occurrences.
const ProxyError = "unable to connect to target host"

// TooManyProxyError is the error returned instead of ProxyError, when it occurs
// too many times.
const TooManyProxyError = "too many proxy failures to target host"
