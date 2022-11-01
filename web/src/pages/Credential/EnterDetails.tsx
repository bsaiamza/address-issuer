import { useEffect, useState } from 'react'
import {
	Alert,
	Box,
	Button,
	FormControl,
	Grid,
	InputLabel,
	MenuItem,
	Paper,
	Select,
	TextField,
	Typography,
	useTheme,
} from '@mui/material'
import { Formik, Form } from 'formik'
import { toast } from 'react-toastify'
import axios from 'axios'
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'
import { DatePicker } from '@mui/x-date-pickers/DatePicker'
import { FileUploader } from 'react-drag-drop-files'
import { AddAPhoto } from '@mui/icons-material'

import { ADDRESS_ISSUER_URL, GMAK, LIGHT_MODE_THEME } from '../../utils/constants'

interface EnterDetailsProps {
	handleNext: () => void
}

const apiURL = ADDRESS_ISSUER_URL + '/credential'
const apiEmailURL = ADDRESS_ISSUER_URL + '/email-credential'

const fileTypes = ['JPEG', 'JPG']

const EnterDetails = ({ handleNext }: EnterDetailsProps) => {
	const theme = useTheme()

	const [submitting, setSubmitting] = useState<boolean | undefined>(false)
	const [input, setInput] = useState('')
	const [streetNumber, setStreetNumber] = useState('')
	const [streetName, setStreetName] = useState('')
	const [area, setArea] = useState('')
	const [locality, setLocality] = useState('')
	const [city, setCity] = useState('')
	const [province, setProvince] = useState('')
	const [postalCode, setPostalCode] = useState('')
	const [countryCode, setCountryCode] = useState('')
	const [file, setFile] = useState(null)

	const inputId = 'address'

	useEffect(() => {
		const renderGoogle = () => {
			// @ts-ignore
			window[inputId] = new window.google.maps.places.Autocomplete(document.getElementById(inputId), {})
			const handlePlaceSelect = () => {
				// @ts-ignore
				const place = window[inputId].getPlace()
				for (const component of place.address_components) {
					const type = component.types[0]
					switch (type) {
						case 'street_number':
							setStreetNumber(component.long_name)
							break
						case 'route':
							setStreetName(component.long_name)
							break
						case 'sublocality_level_1':
							setArea(component.long_name)
							break
						case 'locality':
							setLocality(component.long_name)
							break
						case 'administrative_area_level_2':
							setCity(component.long_name)
							break
						case 'administrative_area_level_1':
							setProvince(component.short_name)
							break
						case 'postal_code':
							setPostalCode(component.long_name)
							break
						case 'country':
							setCountryCode(component.short_name)
							break

						default:
							console.log('irrelevant component type')
							break
					}
				}
			}

			//listen for place change in input field
			// @ts-ignore
			window[inputId].addListener('place_changed', handlePlaceSelect)
		}

		//if places script is already found then listen for load and then renderGoogle
		let found = document.getElementById('placesScript') ? true : false
		if (!found) {
			const script = document.createElement('script')
			script.id = 'placesScript'
			script.src = 'https://maps.googleapis.com/maps/api/js?key=' + GMAK + '&libraries=places'
			script.async = true
			script.onload = () => renderGoogle()
			document.body.appendChild(script)
		}
		if (found) {
			// @ts-ignore
			document.getElementById('placesScript').addEventListener('load', renderGoogle)
		}
	}, [inputId, city, streetNumber, province])

	const sendOffer = async (values: any) => {
		setSubmitting(true)

		await toast.promise(
			axios
				.post(apiURL, values)
				.then((response: any) => {
					// @ts-ignore
					sessionStorage.setItem('credential', response.data.credential)
					toast.success('Credential generated!')
					setTimeout(() => {
						handleNext()
					}, 1000)
				})
				.catch((error: any) => {
					toast.error(error.response.data.msg)
				}),
			{
				pending: 'Generating credential...',
			}
		)
		setSubmitting(false)
	}

	const emailOffer = async (values: any) => {
		setSubmitting(true)

		await toast.promise(
			axios
				.post(apiEmailURL, values)
				.then((response: any) => {
					// @ts-ignore
					sessionStorage.setItem('credential', response.data.credential)
					toast.success('Emailed credential!')
				})
				.catch((error: any) => {
					toast.error(error.response.data.msg)
				}),
			{
				pending: 'Emailing credential...',
			}
		)
		setSubmitting(false)
	}

	const handleFileChange = (file: any) => {
		getBase64(file, (result: any) => {
			const imgBase64 = result
			setFile(imgBase64)
		})
	}

	const getBase64 = (file: any, cb: any) => {
		let reader = new FileReader()
		reader.readAsDataURL(file)
		reader.onload = function () {
			cb(reader.result)
		}
		reader.onerror = function (error) {
			console.log('Error: ', error)
		}
	}

	return (
		<Grid container>
			<Grid item xs={0} md={2} />
			<Grid item xs={12} md={8}>
				<Paper
					square
					elevation={2}
					sx={{ p: 5, width: { md: '100%' }, backgroundColor: theme.palette.mode === LIGHT_MODE_THEME ? '#fff' : '', borderRadius: 5 }}
					>
					<Formik
						initialValues={{
							address_line: '',
							city: '',
							province: '',
							postal_code: '',
							country_code: '',
							statement_date: '',
							statement_image: '',
							self_attested: '',
							email: '',
						}}
						// validate={idValidation}
						onSubmit={(values) => {
							values.email === '' ? sendOffer(values) : emailOffer(values)
						}}>
						{({ values, handleChange, setFieldValue }) => (
							<Form>
								<div>
									<TextField
										id={inputId}
										name='search'
										label='Search for your address'
										value={input}
										onChange={(e) => setInput(e.target.value)}
										sx={{ m: '1rem' }}
										required
										fullWidth
									/>
								</div>
								<div>
									<TextField
										id='address_line'
										name='address_line'
										value={(values.address_line = streetNumber + ' ' + streetName + ' ' + area + ' ' + locality)}
										onChange={handleChange}
										label='Address'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										disabled
									/>
								</div>

								<div>
									<TextField
										id='city'
										name='city'
										value={(values.city = city)}
										onChange={handleChange}
										label='City'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										disabled
									/>
								</div>

								<div>
									<TextField
										id='province'
										name='province'
										value={(values.province = province)}
										onChange={handleChange}
										label='Province'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										disabled
									/>
								</div>

								<div>
									<TextField
										id='postal_code'
										name='postal_code'
										value={(values.postal_code = postalCode)}
										onChange={handleChange}
										label='Postal Code'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										disabled
									/>
								</div>

								<div>
									<TextField
										id='country_code'
										name='country_code'
										value={(values.country_code = countryCode)}
										onChange={handleChange}
										label='Country Code'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										disabled
									/>
								</div>

								<div>
									<LocalizationProvider dateAdapter={AdapterDayjs}>
										<DatePicker
											label='Statement Date'
											value={values.statement_date}
											onChange={(value) => setFieldValue('statement_date', value, true)}
											renderInput={(params) => (
												<TextField
													id='statement_date'
													name='statement_date'
													{...params}
													sx={{ m: '1rem', width: '14.5rem' }}
													required
												/>
											)}
										/>
									</LocalizationProvider>

									<FormControl sx={{ width: '16.5rem' }}>
										<InputLabel id='self_attested' sx={{ margin: '1rem 0 0 1rem' }} required>
											Self Attested
										</InputLabel>
										<Select
											labelId='self_attested'
											id='self_attested'
											name='self_attested'
											value={values.self_attested}
											label='Self Attested'
											onChange={handleChange}
											required
											sx={{ m: '1rem' }}>
											<MenuItem value='N'>N</MenuItem>
											<MenuItem value='Y'>Y</MenuItem>
										</Select>
									</FormControl>
								</div>

								<FileUploader
									name='statement'
									handleChange={handleFileChange}
									types={fileTypes}
									maxSize={0.1}
									onSizeError={() => toast.error('Image is too large. Max size is 100kb.')}
									onTypeError={() => toast.error('Incorrect file type.')}>
									<Box
										sx={{
											width: '100%',
											height: '20rem',
											borderRadius: '0.5rem',
											border: '1px dashed #de5353',
											p: file ? '1rem' : '5rem',
											cursor: 'pointer',
										}}>
										{file ? (
											<>
												<img src={file} alt='statement' style={{ maxWidth: 400, maxHeight: 280 }} />
												<TextField
													id='statement_image'
													name='statement_image'
													value={(values.statement_image = file)}
													required
													sx={{ display: 'none' }}
												/>
											</>
										) : (
											<>
												<Typography variant='h5'>
													<AddAPhoto /> Upload or Drag an image of your statement
												</Typography>
												<Typography>Type: JPEG or JPG | Max Size: 100kb</Typography>
											</>
										)}
									</Box>
								</FileUploader>

								<div>
									<TextField
										id='email'
										name='email'
										type='email'
										value={values.email}
										onChange={handleChange}
										label='Email'
										sx={{ m: '1rem' }}
										// helperText='Optional: If you want to receive your credential via email.'
										helperText={
											<Alert severity='info' sx={{ backgroundColor: 'transparent' }}>
												Optional: If you want to receive your credential via email.
											</Alert>
										}
									/>
								</div>

								<div>
									<Button
										variant='contained'
										size='small'
										type='submit'
										sx={{ color: '#fff', m: '1rem' }}
										disabled={submitting}>
										Submit
									</Button>
								</div>
							</Form>
						)}
					</Formik>
				</Paper>
			</Grid>
			<Grid item xs={0} md={2} />
		</Grid>
	)
}

export default EnterDetails
